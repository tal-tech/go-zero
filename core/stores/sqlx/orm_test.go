package sqlx

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
)

func TestUnmarshalRowBool(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value bool
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.True(t, value)
	})
}

func TestUnmarshalRowBoolNotSettable(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value bool
		assert.NotNil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select value from users where user=?", "anyone"))
	})
}

func TestUnmarshalRowInt(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value int
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, 2, value)
	})
}

func TestUnmarshalRowInt8(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value int8
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, int8(3), value)
	})
}

func TestUnmarshalRowInt16(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("4")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value int16
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.Equal(t, int16(4), value)
	})
}

func TestUnmarshalRowInt32(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("5")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value int32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.Equal(t, int32(5), value)
	})
}

func TestUnmarshalRowInt64(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("6")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value int64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, int64(6), value)
	})
}

func TestUnmarshalRowUint(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value uint
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint(2), value)
	})
}

func TestUnmarshalRowUint8(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value uint8
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint8(3), value)
	})
}

func TestUnmarshalRowUint16(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("4")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value uint16
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint16(4), value)
	})
}

func TestUnmarshalRowUint32(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("5")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value uint32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint32(5), value)
	})
}

func TestUnmarshalRowUint64(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("6")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value uint64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint16(6), value)
	})
}

func TestUnmarshalRowFloat32(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("7")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value float32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, float32(7), value)
	})
}

func TestUnmarshalRowFloat64(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("8")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value float64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, float64(8), value)
	})
}

func TestUnmarshalRowString(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		const expect = "hello"
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString(expect)
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value string
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowStruct(t *testing.T) {
	var value = new(struct {
		Name string `db:"name"`
		Age  int    `db:"age"`
	})

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, "liao", value.Name)
		assert.Equal(t, 5, value.Age)
	})
}

func TestUnmarshalRowStructWithTags(t *testing.T) {
	var value = new(struct {
		Age  int    `db:"age"`
		Name string `db:"name"`
	})

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, "liao", value.Name)
		assert.Equal(t, 5, value.Age)
	})
}

func TestUnmarshalRowStructWithTagsWrongColumns(t *testing.T) {
	var value = new(struct {
		Age  *int   `db:"age"`
		Name string `db:"name"`
	})

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name"}).FromCSVString("liao")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		assert.NotNil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
	})
}

func TestUnmarshalRowsBool(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []bool{true, false}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []bool
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []int{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []int
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt8(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []int8{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []int8
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt16(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []int16{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []int16
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt32(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []int32{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []int32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt64(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []int64{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []int64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []uint{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []uint
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint8(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []uint8{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []uint8
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint16(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []uint16{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []uint16
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint32(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []uint32{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []uint32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint64(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []uint64{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []uint64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat32(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []float32{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []float32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat64(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []float64{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []float64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsString(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []string{"hello", "world"}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("hello\nworld")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []string
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsBoolPtr(t *testing.T) {
	yes := true
	no := false
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*bool{&yes, &no}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*bool
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsIntPtr(t *testing.T) {
	two := 2
	three := 3
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*int{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*int
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt8Ptr(t *testing.T) {
	two := int8(2)
	three := int8(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*int8{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*int8
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt16Ptr(t *testing.T) {
	two := int16(2)
	three := int16(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*int16{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*int16
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt32Ptr(t *testing.T) {
	two := int32(2)
	three := int32(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*int32{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*int32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt64Ptr(t *testing.T) {
	two := int64(2)
	three := int64(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*int64{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*int64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUintPtr(t *testing.T) {
	two := uint(2)
	three := uint(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*uint{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*uint
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint8Ptr(t *testing.T) {
	two := uint8(2)
	three := uint8(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*uint8{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*uint8
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint16Ptr(t *testing.T) {
	two := uint16(2)
	three := uint16(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*uint16{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*uint16
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint32Ptr(t *testing.T) {
	two := uint32(2)
	three := uint32(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*uint32{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*uint32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint64Ptr(t *testing.T) {
	two := uint64(2)
	three := uint64(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*uint64{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*uint64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat32Ptr(t *testing.T) {
	two := float32(2)
	three := float32(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*float32{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*float32
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat64Ptr(t *testing.T) {
	two := float64(2)
	three := float64(3)
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*float64{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*float64
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsStringPtr(t *testing.T) {
	hello := "hello"
	world := "world"
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var expect = []*string{&hello, &world}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("hello\nworld")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var value []*string
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsStruct(t *testing.T) {
	var expect = []struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []struct {
		Name string
		Age  int64
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructWithNullStringType(t *testing.T) {
	var expect = []struct {
		Name       string
		NullString sql.NullString
	}{
		{
			Name: "first",
			NullString: sql.NullString{
				String: "firstnullstring",
				Valid:  true,
			},
		},
		{
			Name: "second",
			NullString: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	var value []struct {
		Name       string         `db:"name"`
		NullString sql.NullString `db:"value"`
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "value"}).AddRow(
			"first", "firstnullstring").AddRow("second", nil)
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.NullString.String, value[i].NullString.String)
			assert.Equal(t, each.NullString.Valid, value[i].NullString.Valid)
		}
	})
}

func TestUnmarshalRowsStructWithTags(t *testing.T) {
	var expect = []struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []struct {
		Age  int64  `db:"age"`
		Name string `db:"name"`
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructAndEmbeddedAnonymousStructWithTags(t *testing.T) {
	type Embed struct {
		Value int64 `db:"value"`
	}

	var expect = []struct {
		Name  string
		Age   int64
		Value int64
	}{
		{
			Name:  "first",
			Age:   2,
			Value: 3,
		},
		{
			Name:  "second",
			Age:   3,
			Value: 4,
		},
	}
	var value []struct {
		Name string `db:"name"`
		Age  int64  `db:"age"`
		Embed
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age", "value"}).FromCSVString("first,2,3\nsecond,3,4")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age, value from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.Value, value[i].Value)
		}
	})
}

func TestUnmarshalRowsStructAndEmbeddedStructPtrAnonymousWithTags(t *testing.T) {
	type Embed struct {
		Value int64 `db:"value"`
	}

	var expect = []struct {
		Name  string
		Age   int64
		Value int64
	}{
		{
			Name:  "first",
			Age:   2,
			Value: 3,
		},
		{
			Name:  "second",
			Age:   3,
			Value: 4,
		},
	}
	var value []struct {
		Name string `db:"name"`
		Age  int64  `db:"age"`
		*Embed
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age", "value"}).FromCSVString("first,2,3\nsecond,3,4")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age, value from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.Value, value[i].Value)
		}
	})
}

func TestUnmarshalRowsStructPtr(t *testing.T) {
	var expect = []*struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []*struct {
		Name string
		Age  int64
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructWithTagsPtr(t *testing.T) {
	var expect = []*struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []*struct {
		Age  int64  `db:"age"`
		Name string `db:"name"`
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructWithTagsPtrWithInnerPtr(t *testing.T) {
	var expect = []*struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []*struct {
		Age  *int64 `db:"age"`
		Name string `db:"name"`
	}

	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, *value[i].Age)
		}
	})
}

func TestCommonSqlConn_QueryRowOptional(t *testing.T) {
	runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"age"}).FromCSVString("5")
		mock.ExpectQuery("select (.+) from users where user=?").WithArgs("anyone").WillReturnRows(rs)

		var r struct {
			User string `db:"user"`
			Age  int    `db:"age"`
		}
		assert.Nil(t, query(db, func(rows *sql.Rows) error {
			return unmarshalRow(&r, rows, false)
		}, "select age from users where user=?", "anyone"))
		assert.Empty(t, r.User)
		assert.Equal(t, 5, r.Age)
	})
}

func TestUnmarshalRowError(t *testing.T) {
	tests := []struct {
		name     string
		colErr   error
		scanErr  error
		err      error
		next     int
		validate func(err error)
	}{
		{
			name: "with error",
			err:  errors.New("foo"),
			validate: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "without next",
			validate: func(err error) {
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:    "with error",
			scanErr: errors.New("foo"),
			next:    1,
			validate: func(err error) {
				assert.Equal(t, ErrNotFound, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runOrmTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"age"}).FromCSVString("5")
				mock.ExpectQuery("select (.+) from users where user=?").WithArgs(
					"anyone").WillReturnRows(rs)

				var r struct {
					User string `db:"user"`
					Age  int    `db:"age"`
				}
				test.validate(query(db, func(rows *sql.Rows) error {
					scanner := mockedScanner{
						colErr:  test.colErr,
						scanErr: test.scanErr,
						err:     test.err,
					}
					return unmarshalRow(&r, &scanner, false)
				}, "select age from users where user=?", "anyone"))
			})
		})
	}
}

func runOrmTest(t *testing.T, fn func(db *sql.DB, mock sqlmock.Sqlmock)) {
	logx.Disable()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fn(db, mock)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

type mockedScanner struct {
	colErr  error
	scanErr error
	err     error
	next    int
}

func (m *mockedScanner) Columns() ([]string, error) {
	return nil, m.colErr
}

func (m *mockedScanner) Err() error {
	return m.err
}

func (m *mockedScanner) Next() bool {
	if m.next > 0 {
		m.next--
		return true
	}
	return false
}

func (m *mockedScanner) Scan(v ...interface{}) error {
	return m.scanErr
}
