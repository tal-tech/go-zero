package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	modelutil "github.com/tal-tech/go-zero/tools/goctl/model/sql/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	pwd             = "."
	createTableFlag = `(?m)^(?i)CREATE\s+TABLE` // ignore case
)

type (
	defaultGenerator struct {
		// source string
		dir string
		console.Console
		pkg string
		cfg *config.Config
	}

	// Option defines a function with argument defaultGenerator
	Option func(generator *defaultGenerator)

	code struct {
		importsCode string
		varsCode    string
		typesCode   string
		newCode     string
		insertCode  string
		findCode    []string
		updateCode  string
		deleteCode  string
		cacheExtra  string
	}
)

// NewDefaultGenerator creates an instance for defaultGenerator
func NewDefaultGenerator(dir string, cfg *config.Config, opt ...Option) (*defaultGenerator, error) {
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

	generator := &defaultGenerator{dir: dir, cfg: cfg, pkg: pkg}
	var optionList []Option
	optionList = append(optionList, newDefaultOption())
	optionList = append(optionList, opt...)
	for _, fn := range optionList {
		fn(generator)
	}

	return generator, nil
}

// WithConsoleOption creates a console option
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

func (g *defaultGenerator) StartFromDDL(filename string, withCache bool) error {
	modelList, err := g.genFromDDL(filename, withCache)
	if err != nil {
		return err
	}

	return g.createFile(modelList)
}

func (g *defaultGenerator) StartFromInformationSchema(tables map[string]*model.Table, withCache bool) error {
	m := make(map[string]string)
	for _, each := range tables {
		table, err := parser.ConvertDataType(each)
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
		modelFilename, err := format.FileNamingFormat(g.cfg.NamingFormat, fmt.Sprintf("%s_model", tn.Source()))
		if err != nil {
			return err
		}

		name := modelFilename + ".go"
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
	varFilename, err := format.FileNamingFormat(g.cfg.NamingFormat, "vars")
	if err != nil {
		return err
	}

	filename := filepath.Join(dirAbs, varFilename+".go")
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
func (g *defaultGenerator) genFromDDL(filename string, withCache bool) (map[string]string, error) {
	m := make(map[string]string)
	tables, err := parser.Parse(filename)
	if err != nil {
		return nil, err
	}

	for _, e := range tables {
		code, err := g.genModel(*e, withCache)
		if err != nil {
			return nil, err
		}

		m[e.Name.Source()] = code
	}

	return m, nil
}

// Table defines mysql table
type Table struct {
	parser.Table
	PrimaryCacheKey        Key
	UniqueCacheKey         []Key
	ContainsUniqueCacheKey bool
}

func (g *defaultGenerator) genModel(in parser.Table, withCache bool) (string, error) {
	if len(in.PrimaryKey.Name.Source()) == 0 {
		return "", fmt.Errorf("table %s: missing primary key", in.Name.Source())
	}

	primaryKey, uniqueKey := genCacheKeys(in)

	importsCode, err := genImports(withCache, in.ContainsTime())
	if err != nil {
		return "", err
	}

	var table Table
	table.Table = in
	table.PrimaryCacheKey = primaryKey
	table.UniqueCacheKey = uniqueKey
	table.ContainsUniqueCacheKey = len(uniqueKey) > 0

	varsCode, err := genVars(table, withCache)
	if err != nil {
		return "", err
	}

	insertCode, insertCodeMethod, err := genInsert(table, withCache)
	if err != nil {
		return "", err
	}

	findCode := make([]string, 0)
	findOneCode, findOneCodeMethod, err := genFindOne(table, withCache)
	if err != nil {
		return "", err
	}

	ret, err := genFindOneByField(table, withCache)
	if err != nil {
		return "", err
	}

	findCode = append(findCode, findOneCode, ret.findOneMethod)
	updateCode, updateCodeMethod, err := genUpdate(table, withCache)
	if err != nil {
		return "", err
	}

	deleteCode, deleteCodeMethod, err := genDelete(table, withCache)
	if err != nil {
		return "", err
	}

	var list []string
	list = append(list, insertCodeMethod, findOneCodeMethod, ret.findOneInterfaceMethod, updateCodeMethod, deleteCodeMethod)
	typesCode, err := genTypes(table, strings.Join(modelutil.TrimStringSlice(list), util.NL), withCache)
	if err != nil {
		return "", err
	}

	newCode, err := genNew(table, withCache)
	if err != nil {
		return "", err
	}

	code := &code{
		importsCode: importsCode,
		varsCode:    varsCode,
		typesCode:   typesCode,
		newCode:     newCode,
		insertCode:  insertCode,
		findCode:    findCode,
		updateCode:  updateCode,
		deleteCode:  deleteCode,
		cacheExtra:  ret.cacheExtra,
	}

	output, err := g.executeModel(code)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func (g *defaultGenerator) executeModel(code *code) (*bytes.Buffer, error) {
	text, err := util.LoadTemplate(category, modelTemplateFile, template.Model)
	if err != nil {
		return nil, err
	}
	t := util.With("model").
		Parse(text).
		GoFmt(true)
	output, err := t.Execute(map[string]interface{}{
		"pkg":         g.pkg,
		"imports":     code.importsCode,
		"vars":        code.varsCode,
		"types":       code.typesCode,
		"new":         code.newCode,
		"insert":      code.insertCode,
		"find":        strings.Join(code.findCode, "\n"),
		"update":      code.updateCode,
		"delete":      code.deleteCode,
		"extraMethod": code.cacheExtra,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func wrapWithRawString(v string) string {
	if v == "`" {
		return v
	}

	if !strings.HasPrefix(v, "`") {
		v = "`" + v
	}

	if !strings.HasSuffix(v, "`") {
		v = v + "`"
	} else if len(v) == 1 {
		v = v + "`"
	}

	return v
}
