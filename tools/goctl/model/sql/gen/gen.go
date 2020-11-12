package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	pwd             = "."
	createTableFlag = `(?m)^(?i)CREATE\s+TABLE` // ignore case
	NamingLower     = "lower"
	NamingCamel     = "camel"
	NamingSnake     = "snake"
)

type (
	defaultGenerator struct {
		//source string
		dir string
		console.Console
		pkg         string
		namingStyle string
	}
	Option func(generator *defaultGenerator)
)

func NewDefaultGenerator(dir, namingStyle string, opt ...Option) (*defaultGenerator, error) {
	if dir == "" {
		dir = pwd
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	dir = dirAbs
	pkg := filepath.Base(dirAbs)
	err = util.MkdirIfNotExist(dir)
	if err != nil {
		return nil, err
	}

	generator := &defaultGenerator{dir: dir, namingStyle: namingStyle, pkg: pkg}
	var optionList []Option
	optionList = append(optionList, newDefaultOption())
	optionList = append(optionList, opt...)
	for _, fn := range optionList {
		fn(generator)
	}
	return generator, nil
}

func WithConsoleOption(c console.Console) Option {
	return func(generator *defaultGenerator) {
		generator.Console = c
	}
}

func newDefaultOption() Option {
	return func(generator *defaultGenerator) {
		generator.Console = console.NewColorConsole()
	}
}

func (g *defaultGenerator) StartFromDdl(source string, withCache bool) error {
	modelList, err := g.genFromDDL(source, withCache)
	if err != nil {
		return err
	}
	return g.createFile(modelList)
}

func (g *defaultGenerator) StartFromInformationSchema(db string, columns map[string][]*model.Column, withCache bool) error {
	m := make(map[string]string)
	for tableName, column := range columns {
		table, err := parser.ConvertColumn(db, tableName, column)
		if err != nil {
			return err
		}
		code, err := g.genModel(*table, withCache)
		if err != nil {
			return err
		}
		m[table.Name.Source()] = code
	}
	return g.createFile(m)
}

func (g *defaultGenerator) createFile(modelList map[string]string) error {
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}
	g.dir = dirAbs
	g.pkg = filepath.Base(dirAbs)
	err = util.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}

	for tableName, code := range modelList {
		tn := stringx.From(tableName)
		name := fmt.Sprintf("%smodel.go", strings.ToLower(tn.ToCamel()))
		switch g.namingStyle {
		case NamingCamel:
			name = fmt.Sprintf("%sModel.go", tn.ToCamel())
		case NamingSnake:
			name = fmt.Sprintf("%s_model.go", tn.ToSnake())
		}
		filename := filepath.Join(dirAbs, name)
		if util.FileExists(filename) {
			g.Warning("%s already exists, ignored.", name)
			continue
		}
		err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
	}
	// generate error file
	filename := filepath.Join(dirAbs, "vars.go")
	if g.namingStyle == NamingCamel {
		filename = filepath.Join(dirAbs, "Vars.go")
	}
	text, err := util.LoadTemplate(category, errTemplateFile, template.Error)
	if err != nil {
		return err
	}

	err = util.With("vars").Parse(text).SaveTo(map[string]interface{}{
		"pkg": g.pkg,
	}, filename, false)
	if err != nil {
		return err
	}

	g.Success("Done.")
	return nil
}

// ret1: key-table name,value-code
func (g *defaultGenerator) genFromDDL(source string, withCache bool) (map[string]string, error) {
	ddlList := g.split(source)
	m := make(map[string]string)
	for _, ddl := range ddlList {
		table, err := parser.Parse(ddl)
		if err != nil {
			return nil, err
		}
		code, err := g.genModel(*table, withCache)
		if err != nil {
			return nil, err
		}
		m[table.Name.Source()] = code
	}
	return m, nil
}

type (
	Table struct {
		parser.Table
		CacheKey          map[string]Key
		ContainsUniqueKey bool
	}
)

func (g *defaultGenerator) genModel(in parser.Table, withCache bool) (string, error) {
	if len(in.PrimaryKey.Name.Source()) == 0 {
		return "", fmt.Errorf("table %s: primary key can not be nil", in.Name.Source())
	}
	text, err := util.LoadTemplate(category, modelTemplateFile, template.Model)
	if err != nil {
		return "", err
	}
	t := util.With("model").
		Parse(text).
		GoFmt(true)

	m, err := genCacheKeys(in)
	if err != nil {
		return "", err
	}

	importsCode, err := genImports(withCache, in.ContainsTime())
	if err != nil {
		return "", err
	}

	var table Table
	table.Table = in
	table.CacheKey = m
	var containsUniqueCache = false
	for _, item := range table.Fields {
		if item.IsUniqueKey {
			containsUniqueCache = true
			break
		}
	}
	table.ContainsUniqueKey = containsUniqueCache

	varsCode, err := genVars(table, withCache)
	if err != nil {
		return "", err
	}

	typesCode, err := genTypes(table, withCache)
	if err != nil {
		return "", err
	}

	newCode, err := genNew(table, withCache)
	if err != nil {
		return "", err
	}

	insertCode, err := genInsert(table, withCache)
	if err != nil {
		return "", err
	}

	var findCode = make([]string, 0)
	findOneCode, err := genFindOne(table, withCache)
	if err != nil {
		return "", err
	}

	findOneByFieldCode, extraMethod, err := genFindOneByField(table, withCache)
	if err != nil {
		return "", err
	}

	findCode = append(findCode, findOneCode, findOneByFieldCode)
	updateCode, err := genUpdate(table, withCache)
	if err != nil {
		return "", err
	}

	deleteCode, err := genDelete(table, withCache)
	if err != nil {
		return "", err
	}

	output, err := t.Execute(map[string]interface{}{
		"pkg":         g.pkg,
		"imports":     importsCode,
		"vars":        varsCode,
		"types":       typesCode,
		"new":         newCode,
		"insert":      insertCode,
		"find":        strings.Join(findCode, "\n"),
		"update":      updateCode,
		"delete":      deleteCode,
		"extraMethod": extraMethod,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
