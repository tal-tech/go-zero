package gen

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var source = "CREATE TABLE `test_user` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `mobile` varchar(255) COLLATE utf8mb4_bin NOT NULL,\n  `class` bigint NOT NULL,\n  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `mobile_unique` (`mobile`),\n  UNIQUE KEY `class_name_unique` (`class`,`name`),\n  KEY `create_index` (`create_time`),\n  KEY `name_index` (`name`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;"

func TestCacheModel(t *testing.T) {
	logx.Disable()
	_ = Clean()

	sqlFile := filepath.Join(t.TempDir(), "tmp.sql")
	err := ioutil.WriteFile(sqlFile, []byte(source), 0777)
	assert.Nil(t, err)

	dir := filepath.Join(t.TempDir(), "./testmodel")
	cacheDir := filepath.Join(dir, "cache")
	noCacheDir := filepath.Join(dir, "nocache")
	g, err := NewDefaultGenerator(cacheDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(cacheDir, "TestUserModel.go"))
		return err == nil
	}())
	g, err = NewDefaultGenerator(noCacheDir, &config.Config{
		NamingFormat: "gozero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, false)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(noCacheDir, "testusermodel.go"))
		return err == nil
	}())
}

func TestNamingModel(t *testing.T) {
	logx.Disable()
	_ = Clean()

	sqlFile := filepath.Join(t.TempDir(), "tmp.sql")
	err := ioutil.WriteFile(sqlFile, []byte(source), 0777)
	assert.Nil(t, err)

	dir, _ := filepath.Abs("./testmodel")
	camelDir := filepath.Join(dir, "camel")
	snakeDir := filepath.Join(dir, "snake")
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	g, err := NewDefaultGenerator(camelDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(camelDir, "TestUserModel.go"))
		return err == nil
	}())
	g, err = NewDefaultGenerator(snakeDir, &config.Config{
		NamingFormat: "go_zero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(snakeDir, "test_user_model.go"))
		return err == nil
	}())
}

func TestWrapWithRawString(t *testing.T) {
	assert.Equal(t, "``", wrapWithRawString(""))
	assert.Equal(t, "``", wrapWithRawString("``"))
	assert.Equal(t, "`a`", wrapWithRawString("a"))
	assert.Equal(t, "`   `", wrapWithRawString("   "))
}

func TestFields(t *testing.T) {
	type Student struct {
		ID         int64           `db:"id"`
		Name       string          `db:"name"`
		Age        sql.NullInt64   `db:"age"`
		Score      sql.NullFloat64 `db:"score"`
		CreateTime time.Time       `db:"create_time"`
		UpdateTime sql.NullTime    `db:"update_time"`
	}
	var (
		studentFieldNames          = builderx.RawFieldNames(&Student{})
		studentRows                = strings.Join(studentFieldNames, ",")
		studentRowsExpectAutoSet   = strings.Join(stringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
		studentRowsWithPlaceHolder = strings.Join(stringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	)

	assert.Equal(t, []string{"`id`", "`name`", "`age`", "`score`", "`create_time`", "`update_time`"}, studentFieldNames)
	assert.Equal(t, "`id`,`name`,`age`,`score`,`create_time`,`update_time`", studentRows)
	assert.Equal(t, "`name`,`age`,`score`", studentRowsExpectAutoSet)
	assert.Equal(t, "`name`=?,`age`=?,`score`=?", studentRowsWithPlaceHolder)
}
