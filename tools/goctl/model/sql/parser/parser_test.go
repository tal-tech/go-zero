package parser

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/util"
)

func TestParsePlainText(t *testing.T) {
	sqlFile := filepath.Join(t.TempDir(), "tmp.sql")
	err := ioutil.WriteFile(sqlFile, []byte("plain text"), 0777)
	assert.Nil(t, err)

	_, err = Parse(sqlFile)
	assert.NotNil(t, err)
}

func TestParseSelect(t *testing.T) {
	sqlFile := filepath.Join(t.TempDir(), "tmp.sql")
	err := ioutil.WriteFile(sqlFile, []byte("select * from user"), 0777)
	assert.Nil(t, err)

	tables, err := Parse(sqlFile)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tables))
}

func TestParseCreateTable(t *testing.T) {
	sqlFile := filepath.Join(t.TempDir(), "tmp.sql")
	err := ioutil.WriteFile(sqlFile, []byte("CREATE TABLE `test_user` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `mobile` varchar(255) COLLATE utf8mb4_bin NOT NULL comment '手\\t机  号',\n  `class` bigint NOT NULL comment '班级',\n  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL comment '姓\n  名',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP comment '创建\\r时间',\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `mobile_unique` (`mobile`),\n  UNIQUE KEY `class_name_unique` (`class`,`name`),\n  KEY `create_index` (`create_time`),\n  KEY `name_index` (`name`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;"), 0777)
	assert.Nil(t, err)

	tables, err := Parse(sqlFile)
	assert.Equal(t, 1, len(tables))
	table := tables[0]
	assert.Nil(t, err)
	assert.Equal(t, "test_user", table.Name.Source())
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	assert.Equal(t, true, table.ContainsTime())
	assert.Equal(t, 2, len(table.UniqueIndex))
	assert.True(t, func() bool {
		for _, e := range table.Fields {
			if e.Comment != util.TrimNewLine(e.Comment) {
				return false
			}
		}

		return true
	}())
}

func TestConvertColumn(t *testing.T) {
	t.Run("missingPrimaryKey", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
					},
				},
			},
		}
		_, err := columnData.Convert()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "missing primary key")
	})

	t.Run("jointPrimaryKey", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
					},
					Index: &model.DbIndex{
						IndexName: "PRIMARY",
					},
				},
				{
					DbColumn: &model.DbColumn{
						Name:     "mobile",
						DataType: "varchar",
						Comment:  "手机号",
					},
					Index: &model.DbIndex{
						IndexName: "PRIMARY",
					},
				},
			},
		}
		_, err := columnData.Convert()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "joint primary key is not supported")
	})

	t.Run("normal", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
						Extra:    "auto_increment",
					},
					Index: &model.DbIndex{
						IndexName:  "PRIMARY",
						SeqInIndex: 1,
					},
				},
				{
					DbColumn: &model.DbColumn{
						Name:     "mobile",
						DataType: "varchar",
						Comment:  "手机号",
					},
					Index: &model.DbIndex{
						IndexName:  "mobile_unique",
						SeqInIndex: 1,
					},
				},
			},
		}

		table, err := columnData.Convert()
		assert.Nil(t, err)
		assert.True(t, table.PrimaryKey.Index.IndexName == "PRIMARY" && table.PrimaryKey.Name == "id")
		for _, item := range table.Columns {
			if item.Name == "mobile" {
				assert.True(t, item.Index.NonUnique == 0)
				break
			}
		}
	})
}
