package mapping

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
)

// because json.Number doesn't support strconv.ParseUint(...),
// so we only can test to 62 bits.
const maxUintBitsToTest = 62

func TestUnmarshalWithFullNameNotStruct(t *testing.T) {
	var s map[string]interface{}
	content := []byte(`{"name":"xiaoming"}`)
	err := UnmarshalJsonBytes(content, &s)
	assert.Equal(t, errValueNotStruct, err)
}

func TestUnmarshalWithoutTagName(t *testing.T) {
	type inner struct {
		Optional bool `key:",optional"`
	}
	m := map[string]interface{}{
		"Optional": true,
	}

	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.True(t, in.Optional)
}

func TestUnmarshalBool(t *testing.T) {
	type inner struct {
		True           bool `key:"yes"`
		False          bool `key:"no"`
		TrueFromOne    bool `key:"yesone,string"`
		FalseFromZero  bool `key:"nozero,string"`
		TrueFromTrue   bool `key:"yestrue,string"`
		FalseFromFalse bool `key:"nofalse,string"`
		DefaultTrue    bool `key:"defaulttrue,default=1"`
		Optional       bool `key:"optional,optional"`
	}
	m := map[string]interface{}{
		"yes":     true,
		"no":      false,
		"yesone":  "1",
		"nozero":  "0",
		"yestrue": "true",
		"nofalse": "false",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.True(in.True)
	ast.False(in.False)
	ast.True(in.TrueFromOne)
	ast.False(in.FalseFromZero)
	ast.True(in.TrueFromTrue)
	ast.False(in.FalseFromFalse)
	ast.True(in.DefaultTrue)
}

func TestUnmarshalDuration(t *testing.T) {
	type inner struct {
		Duration     time.Duration `key:"duration"`
		LessDuration time.Duration `key:"less"`
		MoreDuration time.Duration `key:"more"`
	}
	m := map[string]interface{}{
		"duration": "5s",
		"less":     "100ms",
		"more":     "24h",
	}
	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, time.Second*5, in.Duration)
	assert.Equal(t, time.Millisecond*100, in.LessDuration)
	assert.Equal(t, time.Hour*24, in.MoreDuration)
}

func TestUnmarshalDurationDefault(t *testing.T) {
	type inner struct {
		Int      int           `key:"int"`
		Duration time.Duration `key:"duration,default=5s"`
	}
	m := map[string]interface{}{
		"int": 5,
	}
	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, 5, in.Int)
	assert.Equal(t, time.Second*5, in.Duration)
}

func TestUnmarshalDurationPtr(t *testing.T) {
	type inner struct {
		Duration *time.Duration `key:"duration"`
	}
	m := map[string]interface{}{
		"duration": "5s",
	}
	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, time.Second*5, *in.Duration)
}

func TestUnmarshalDurationPtrDefault(t *testing.T) {
	type inner struct {
		Int      int            `key:"int"`
		Value    *int           `key:",default=5"`
		Duration *time.Duration `key:"duration,default=5s"`
	}
	m := map[string]interface{}{
		"int": 5,
	}
	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, 5, in.Int)
	assert.Equal(t, 5, *in.Value)
	assert.Equal(t, time.Second*5, *in.Duration)
}

func TestUnmarshalInt(t *testing.T) {
	type inner struct {
		Int          int   `key:"int"`
		IntFromStr   int   `key:"intstr,string"`
		Int8         int8  `key:"int8"`
		Int8FromStr  int8  `key:"int8str,string"`
		Int16        int16 `key:"int16"`
		Int16FromStr int16 `key:"int16str,string"`
		Int32        int32 `key:"int32"`
		Int32FromStr int32 `key:"int32str,string"`
		Int64        int64 `key:"int64"`
		Int64FromStr int64 `key:"int64str,string"`
		DefaultInt   int64 `key:"defaultint,default=11"`
		Optional     int   `key:"optional,optional"`
	}
	m := map[string]interface{}{
		"int":      1,
		"intstr":   "2",
		"int8":     int8(3),
		"int8str":  "4",
		"int16":    int16(5),
		"int16str": "6",
		"int32":    int32(7),
		"int32str": "8",
		"int64":    int64(9),
		"int64str": "10",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(1, in.Int)
	ast.Equal(2, in.IntFromStr)
	ast.Equal(int8(3), in.Int8)
	ast.Equal(int8(4), in.Int8FromStr)
	ast.Equal(int16(5), in.Int16)
	ast.Equal(int16(6), in.Int16FromStr)
	ast.Equal(int32(7), in.Int32)
	ast.Equal(int32(8), in.Int32FromStr)
	ast.Equal(int64(9), in.Int64)
	ast.Equal(int64(10), in.Int64FromStr)
	ast.Equal(int64(11), in.DefaultInt)
}

func TestUnmarshalIntPtr(t *testing.T) {
	type inner struct {
		Int *int `key:"int"`
	}
	m := map[string]interface{}{
		"int": 1,
	}

	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.NotNil(t, in.Int)
	assert.Equal(t, 1, *in.Int)
}

func TestUnmarshalIntWithDefault(t *testing.T) {
	type inner struct {
		Int int `key:"int,default=5"`
	}
	m := map[string]interface{}{
		"int": 1,
	}

	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, 1, in.Int)
}

func TestUnmarshalUint(t *testing.T) {
	type inner struct {
		Uint          uint   `key:"uint"`
		UintFromStr   uint   `key:"uintstr,string"`
		Uint8         uint8  `key:"uint8"`
		Uint8FromStr  uint8  `key:"uint8str,string"`
		Uint16        uint16 `key:"uint16"`
		Uint16FromStr uint16 `key:"uint16str,string"`
		Uint32        uint32 `key:"uint32"`
		Uint32FromStr uint32 `key:"uint32str,string"`
		Uint64        uint64 `key:"uint64"`
		Uint64FromStr uint64 `key:"uint64str,string"`
		DefaultUint   uint   `key:"defaultuint,default=11"`
		Optional      uint   `key:"optional,optional"`
	}
	m := map[string]interface{}{
		"uint":      uint(1),
		"uintstr":   "2",
		"uint8":     uint8(3),
		"uint8str":  "4",
		"uint16":    uint16(5),
		"uint16str": "6",
		"uint32":    uint32(7),
		"uint32str": "8",
		"uint64":    uint64(9),
		"uint64str": "10",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(uint(1), in.Uint)
	ast.Equal(uint(2), in.UintFromStr)
	ast.Equal(uint8(3), in.Uint8)
	ast.Equal(uint8(4), in.Uint8FromStr)
	ast.Equal(uint16(5), in.Uint16)
	ast.Equal(uint16(6), in.Uint16FromStr)
	ast.Equal(uint32(7), in.Uint32)
	ast.Equal(uint32(8), in.Uint32FromStr)
	ast.Equal(uint64(9), in.Uint64)
	ast.Equal(uint64(10), in.Uint64FromStr)
	ast.Equal(uint(11), in.DefaultUint)
}

func TestUnmarshalFloat(t *testing.T) {
	type inner struct {
		Float32      float32 `key:"float32"`
		Float32Str   float32 `key:"float32str,string"`
		Float64      float64 `key:"float64"`
		Float64Str   float64 `key:"float64str,string"`
		DefaultFloat float32 `key:"defaultfloat,default=5.5"`
		Optional     float32 `key:",optional"`
	}
	m := map[string]interface{}{
		"float32":    float32(1.5),
		"float32str": "2.5",
		"float64":    float64(3.5),
		"float64str": "4.5",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(float32(1.5), in.Float32)
	ast.Equal(float32(2.5), in.Float32Str)
	ast.Equal(3.5, in.Float64)
	ast.Equal(4.5, in.Float64Str)
	ast.Equal(float32(5.5), in.DefaultFloat)
}

func TestUnmarshalInt64Slice(t *testing.T) {
	var v struct {
		Ages []int64 `key:"ages"`
	}
	m := map[string]interface{}{
		"ages": []int64{1, 2},
	}

	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &v))
	ast.ElementsMatch([]int64{1, 2}, v.Ages)
}

func TestUnmarshalIntSlice(t *testing.T) {
	var v struct {
		Ages []int `key:"ages"`
	}
	m := map[string]interface{}{
		"ages": []int{1, 2},
	}

	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &v))
	ast.ElementsMatch([]int{1, 2}, v.Ages)
}

func TestUnmarshalString(t *testing.T) {
	type inner struct {
		Name              string `key:"name"`
		NameStr           string `key:"namestr,string"`
		NotPresent        string `key:",optional"`
		NotPresentWithTag string `key:"notpresent,optional"`
		DefaultString     string `key:"defaultstring,default=hello"`
		Optional          string `key:",optional"`
	}
	m := map[string]interface{}{
		"name":    "kevin",
		"namestr": "namewithstring",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal("kevin", in.Name)
	ast.Equal("namewithstring", in.NameStr)
	ast.Empty(in.NotPresent)
	ast.Empty(in.NotPresentWithTag)
	ast.Equal("hello", in.DefaultString)
}

func TestUnmarshalStringWithMissing(t *testing.T) {
	type inner struct {
		Name string `key:"name"`
	}
	m := map[string]interface{}{}

	var in inner
	assert.NotNil(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStringSliceFromString(t *testing.T) {
	var v struct {
		Names []string `key:"names"`
	}
	m := map[string]interface{}{
		"names": `["first", "second"]`,
	}

	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &v))
	ast.Equal(2, len(v.Names))
	ast.Equal("first", v.Names[0])
	ast.Equal("second", v.Names[1])
}

func TestUnmarshalIntSliceFromString(t *testing.T) {
	var v struct {
		Values []int `key:"values"`
	}
	m := map[string]interface{}{
		"values": `[1, 2]`,
	}

	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &v))
	ast.Equal(2, len(v.Values))
	ast.Equal(1, v.Values[0])
	ast.Equal(2, v.Values[1])
}

func TestUnmarshalStruct(t *testing.T) {
	type address struct {
		City          string `key:"city"`
		ZipCode       int    `key:"zipcode,string"`
		DefaultString string `key:"defaultstring,default=hello"`
		Optional      string `key:",optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}
	m := map[string]interface{}{
		"name": "kevin",
		"address": map[string]interface{}{
			"city":    "shanghai",
			"zipcode": "200000",
		},
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal("kevin", in.Name)
	ast.Equal("shanghai", in.Address.City)
	ast.Equal(200000, in.Address.ZipCode)
	ast.Equal("hello", in.Address.DefaultString)
}

func TestUnmarshalStructOptionalDepends(t *testing.T) {
	type address struct {
		City            string `key:"city"`
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=Optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional":        "a",
				"OptionalDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]interface{}{
				"name": "kevin",
				"address": map[string]interface{}{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]interface{})[k] = v
			}

			var in inner
			ast := assert.New(t)
			if test.pass {
				ast.Nil(UnmarshalKey(m, &in))
				ast.Equal("kevin", in.Name)
				ast.Equal("shanghai", in.Address.City)
				ast.Equal(test.input["Optional"], in.Address.Optional)
				ast.Equal(test.input["OptionalDepends"], in.Address.OptionalDepends)
			} else {
				ast.NotNil(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestUnmarshalStructOptionalDependsNot(t *testing.T) {
	type address struct {
		City            string `key:"city"`
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!Optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			input: map[string]string{},
			pass:  false,
		},
		{
			input: map[string]string{
				"Optional":        "a",
				"OptionalDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional": "a",
			},
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]interface{}{
				"name": "kevin",
				"address": map[string]interface{}{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]interface{})[k] = v
			}

			var in inner
			ast := assert.New(t)
			if test.pass {
				ast.Nil(UnmarshalKey(m, &in))
				ast.Equal("kevin", in.Name)
				ast.Equal("shanghai", in.Address.City)
				ast.Equal(test.input["Optional"], in.Address.Optional)
				ast.Equal(test.input["OptionalDepends"], in.Address.OptionalDepends)
			} else {
				ast.NotNil(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestUnmarshalStructOptionalDependsNotErrorDetails(t *testing.T) {
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!Optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]interface{}{
		"name": "kevin",
	}

	var in inner
	err := UnmarshalKey(m, &in)
	assert.NotNil(t, err)
}

func TestUnmarshalStructOptionalDependsNotNested(t *testing.T) {
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!Optional"`
	}
	type combo struct {
		Name    string  `key:"name,optional"`
		Address address `key:"address"`
	}
	type inner struct {
		Name  string `key:"name"`
		Combo combo  `key:"combo"`
	}

	m := map[string]interface{}{
		"name": "kevin",
	}

	var in inner
	err := UnmarshalKey(m, &in)
	assert.NotNil(t, err)
}

func TestUnmarshalStructOptionalNestedDifferentKey(t *testing.T) {
	type address struct {
		Optional        string `dkey:",optional"`
		OptionalDepends string `key:",optional"`
	}
	type combo struct {
		Name    string  `key:"name,optional"`
		Address address `key:"address"`
	}
	type inner struct {
		Name  string `key:"name"`
		Combo combo  `key:"combo"`
	}

	m := map[string]interface{}{
		"name": "kevin",
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStructOptionalDependsNotEnoughValue(t *testing.T) {
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]interface{}{
		"name":    "kevin",
		"address": map[string]interface{}{},
	}

	var in inner
	err := UnmarshalKey(m, &in)
	assert.NotNil(t, err)
}

func TestUnmarshalAnonymousStructOptionalDepends(t *testing.T) {
	type AnonAddress struct {
		City            string `key:"city"`
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=Optional"`
	}
	type inner struct {
		Name string `key:"name"`
		AnonAddress
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional":        "a",
				"OptionalDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]interface{}{
				"name": "kevin",
				"city": "shanghai",
			}
			for k, v := range test.input {
				m[k] = v
			}

			var in inner
			ast := assert.New(t)
			if test.pass {
				ast.Nil(UnmarshalKey(m, &in))
				ast.Equal("kevin", in.Name)
				ast.Equal("shanghai", in.City)
				ast.Equal(test.input["Optional"], in.Optional)
				ast.Equal(test.input["OptionalDepends"], in.OptionalDepends)
			} else {
				ast.NotNil(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestUnmarshalStructPtr(t *testing.T) {
	type address struct {
		City          string `key:"city"`
		ZipCode       int    `key:"zipcode,string"`
		DefaultString string `key:"defaultstring,default=hello"`
		Optional      string `key:",optional"`
	}
	type inner struct {
		Name    string   `key:"name"`
		Address *address `key:"address"`
	}
	m := map[string]interface{}{
		"name": "kevin",
		"address": map[string]interface{}{
			"city":    "shanghai",
			"zipcode": "200000",
		},
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal("kevin", in.Name)
	ast.Equal("shanghai", in.Address.City)
	ast.Equal(200000, in.Address.ZipCode)
	ast.Equal("hello", in.Address.DefaultString)
}

func TestUnmarshalWithStringIgnored(t *testing.T) {
	type inner struct {
		True    bool    `key:"yes"`
		False   bool    `key:"no"`
		Int     int     `key:"int"`
		Int8    int8    `key:"int8"`
		Int16   int16   `key:"int16"`
		Int32   int32   `key:"int32"`
		Int64   int64   `key:"int64"`
		Uint    uint    `key:"uint"`
		Uint8   uint8   `key:"uint8"`
		Uint16  uint16  `key:"uint16"`
		Uint32  uint32  `key:"uint32"`
		Uint64  uint64  `key:"uint64"`
		Float32 float32 `key:"float32"`
		Float64 float64 `key:"float64"`
	}
	m := map[string]interface{}{
		"yes":     "1",
		"no":      "0",
		"int":     "1",
		"int8":    "3",
		"int16":   "5",
		"int32":   "7",
		"int64":   "9",
		"uint":    "1",
		"uint8":   "3",
		"uint16":  "5",
		"uint32":  "7",
		"uint64":  "9",
		"float32": "1.5",
		"float64": "3.5",
	}

	var in inner
	um := NewUnmarshaler("key", WithStringValues())
	ast := assert.New(t)
	ast.Nil(um.Unmarshal(m, &in))
	ast.True(in.True)
	ast.False(in.False)
	ast.Equal(1, in.Int)
	ast.Equal(int8(3), in.Int8)
	ast.Equal(int16(5), in.Int16)
	ast.Equal(int32(7), in.Int32)
	ast.Equal(int64(9), in.Int64)
	ast.Equal(uint(1), in.Uint)
	ast.Equal(uint8(3), in.Uint8)
	ast.Equal(uint16(5), in.Uint16)
	ast.Equal(uint32(7), in.Uint32)
	ast.Equal(uint64(9), in.Uint64)
	ast.Equal(float32(1.5), in.Float32)
	ast.Equal(3.5, in.Float64)
}

func TestUnmarshalJsonNumberInt64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue int64 = 1 << uint(i)
		strValue := strconv.FormatInt(intValue, 10)
		number := json.Number(strValue)
		m := map[string]interface{}{
			"ID": number,
		}
		var v struct {
			ID int64
		}
		assert.Nil(t, UnmarshalKey(m, &v))
		assert.Equal(t, intValue, v.ID)
	}
}

func TestUnmarshalJsonNumberUint64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue uint64 = 1 << uint(i)
		strValue := strconv.FormatUint(intValue, 10)
		number := json.Number(strValue)
		m := map[string]interface{}{
			"ID": number,
		}
		var v struct {
			ID uint64
		}
		assert.Nil(t, UnmarshalKey(m, &v))
		assert.Equal(t, intValue, v.ID)
	}
}

func TestUnmarshalJsonNumberUint64Ptr(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue uint64 = 1 << uint(i)
		strValue := strconv.FormatUint(intValue, 10)
		number := json.Number(strValue)
		m := map[string]interface{}{
			"ID": number,
		}
		var v struct {
			ID *uint64
		}
		ast := assert.New(t)
		ast.Nil(UnmarshalKey(m, &v))
		ast.NotNil(v.ID)
		ast.Equal(intValue, *v.ID)
	}
}

func TestUnmarshalMapOfInt(t *testing.T) {
	m := map[string]interface{}{
		"Ids": map[string]bool{"first": true},
	}
	var v struct {
		Ids map[string]bool
	}
	assert.Nil(t, UnmarshalKey(m, &v))
	assert.True(t, v.Ids["first"])
}

func TestUnmarshalMapOfStructError(t *testing.T) {
	m := map[string]interface{}{
		"Ids": map[string]interface{}{"first": "second"},
	}
	var v struct {
		Ids map[string]struct {
			Name string
		}
	}
	assert.NotNil(t, UnmarshalKey(m, &v))
}

func TestUnmarshalSlice(t *testing.T) {
	m := map[string]interface{}{
		"Ids": []interface{}{"first", "second"},
	}
	var v struct {
		Ids []string
	}
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &v))
	ast.Equal(2, len(v.Ids))
	ast.Equal("first", v.Ids[0])
	ast.Equal("second", v.Ids[1])
}

func TestUnmarshalSliceOfStruct(t *testing.T) {
	m := map[string]interface{}{
		"Ids": []map[string]interface{}{
			{
				"First":  1,
				"Second": 2,
			},
		},
	}
	var v struct {
		Ids []struct {
			First  int
			Second int
		}
	}
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &v))
	ast.Equal(1, len(v.Ids))
	ast.Equal(1, v.Ids[0].First)
	ast.Equal(2, v.Ids[0].Second)
}

func TestUnmarshalWithStringOptionsCorrect(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]interface{}{
		"value":   "first",
		"correct": "2",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal("first", in.Value)
	ast.Equal("2", in.Correct)
}

func TestUnmarshalStringOptionsWithStringOptionsNotString(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]interface{}{
		"value":   "first",
		"correct": 2,
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	ast := assert.New(t)
	ast.NotNil(unmarshaler.Unmarshal(m, &in))
}

func TestUnmarshalStringOptionsWithStringOptions(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]interface{}{
		"value":   "first",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	ast := assert.New(t)
	ast.Nil(unmarshaler.Unmarshal(m, &in))
	ast.Equal("first", in.Value)
	ast.Equal("2", in.Correct)
}

func TestUnmarshalStringOptionsWithStringOptionsPtr(t *testing.T) {
	type inner struct {
		Value   *string `key:"value,options=first|second"`
		Correct *int    `key:"correct,options=1|2"`
	}
	m := map[string]interface{}{
		"value":   "first",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	ast := assert.New(t)
	ast.Nil(unmarshaler.Unmarshal(m, &in))
	ast.True(*in.Value == "first")
	ast.True(*in.Correct == 2)
}

func TestUnmarshalStringOptionsWithStringOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]interface{}{
		"value":   "third",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	ast := assert.New(t)
	ast.NotNil(unmarshaler.Unmarshal(m, &in))
}

func TestUnmarshalWithStringOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect string `key:"incorrect,options=1|2"`
	}
	m := map[string]interface{}{
		"value":     "first",
		"incorrect": "3",
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(m, &in))
}

func TestUnmarshalWithIntOptionsCorrect(t *testing.T) {
	type inner struct {
		Value  string `key:"value,options=first|second"`
		Number int    `key:"number,options=1|2"`
	}
	m := map[string]interface{}{
		"value":  "first",
		"number": 2,
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal("first", in.Value)
	ast.Equal(2, in.Number)
}

func TestUnmarshalWithIntOptionsCorrectPtr(t *testing.T) {
	type inner struct {
		Value  *string `key:"value,options=first|second"`
		Number *int    `key:"number,options=1|2"`
	}
	m := map[string]interface{}{
		"value":  "first",
		"number": 2,
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.True(*in.Value == "first")
	ast.True(*in.Number == 2)
}

func TestUnmarshalWithIntOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]interface{}{
		"value":     "first",
		"incorrect": 3,
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(m, &in))
}

func TestUnmarshalWithUintOptionsCorrect(t *testing.T) {
	type inner struct {
		Value  string `key:"value,options=first|second"`
		Number uint   `key:"number,options=1|2"`
	}
	m := map[string]interface{}{
		"value":  "first",
		"number": uint(2),
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal("first", in.Value)
	ast.Equal(uint(2), in.Number)
}

func TestUnmarshalWithUintOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect uint   `key:"incorrect,options=1|2"`
	}
	m := map[string]interface{}{
		"value":     "first",
		"incorrect": uint(3),
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(m, &in))
}

func TestUnmarshalWithOptionsAndDefault(t *testing.T) {
	type inner struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]interface{}{}

	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, "second", in.Value)
}

func TestUnmarshalWithOptionsAndSet(t *testing.T) {
	type inner struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]interface{}{
		"value": "first",
	}

	var in inner
	assert.Nil(t, UnmarshalKey(m, &in))
	assert.Equal(t, "first", in.Value)
}

func TestUnmarshalNestedKey(t *testing.T) {
	var c struct {
		ID int `json:"Persons.first.ID"`
	}
	m := map[string]interface{}{
		"Persons": map[string]interface{}{
			"first": map[string]interface{}{
				"ID": 1,
			},
		},
	}

	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &c))
	assert.Equal(t, 1, c.ID)
}

func TestUnmarhsalNestedKeyArray(t *testing.T) {
	var c struct {
		First []struct {
			ID int
		} `json:"Persons.first"`
	}
	m := map[string]interface{}{
		"Persons": map[string]interface{}{
			"first": []map[string]interface{}{
				{"ID": 1},
				{"ID": 2},
			},
		},
	}

	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &c))
	assert.Equal(t, 2, len(c.First))
	assert.Equal(t, 1, c.First[0].ID)
}

func TestUnmarshalAnonymousOptionalRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousOptionalOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousOptionalRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "kevin", b.Name)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOptionalRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "kevin", b.Name)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalRequiredProvidedOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalRequiredMissedOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"n": "anything",
	}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOptionalBothOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "kevin", b.Name)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalBothOptionalOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalBothOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "kevin", b.Name)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOneRequiredOneOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "kevin", b.Name)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOneRequiredOneOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOneRequiredOneOptionalRequiredProvidedOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousOneRequiredOneOptionalRequiredMissedOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"n": "anything",
	}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalAnonymousBothOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "kevin", b.Name)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousBothOptionalOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.Equal(t, "anything", b.Value)
}

func TestUnmarshalAnonymousBothOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.True(t, len(b.Name) == 0)
	assert.True(t, len(b.Value) == 0)
}

func TestUnmarshalAnonymousWrappedToMuch(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]interface{}{
		"Foo": map[string]interface{}{
			"n": "name",
			"v": "anything",
		},
	}

	var b Bar
	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &b))
}

func TestUnmarshalWrappedObject(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Inner Foo
		}
	)
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"v": "anything",
		},
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Inner.Value)
}

func TestUnmarshalWrappedObjectOptional(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Name  string
		}
	)
	m := map[string]interface{}{
		"Name": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Name)
}

func TestUnmarshalWrappedObjectOptionalFilled(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Name  string
		}
	)
	hosts := []string{"1", "2"}
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"Hosts": hosts,
			"Key":   "key",
		},
		"Name": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.EqualValues(t, hosts, b.Inner.Hosts)
	assert.Equal(t, "key", b.Inner.Key)
	assert.Equal(t, "anything", b.Name)
}

func TestUnmarshalWrappedNamedObjectOptional(t *testing.T) {
	type (
		Foo struct {
			Host string
			Key  string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Name  string
		}
	)
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"Ip":  "thehost",
			"Key": "thekey",
		},
		"Name": "anything",
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "thehost", b.Inner.Host)
	assert.Equal(t, "thekey", b.Inner.Key)
	assert.Equal(t, "anything", b.Name)
}

func TestUnmarshalWrappedObjectNamedPtr(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Inner *Foo `json:"foo,optional"`
		}
	)
	m := map[string]interface{}{
		"foo": map[string]interface{}{
			"v": "anything",
		},
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Inner.Value)
}

func TestUnmarshalWrappedObjectPtr(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Inner *Foo
		}
	)
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"v": "anything",
		},
	}

	var b Bar
	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &b))
	assert.Equal(t, "anything", b.Inner.Value)
}

func TestUnmarshalInt2String(t *testing.T) {
	type inner struct {
		Int string `key:"int"`
	}
	m := map[string]interface{}{
		"int": 123,
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(m, &in))
}

func TestUnmarshalZeroValues(t *testing.T) {
	type inner struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `key:"string"`
	}
	m := map[string]interface{}{
		"no":     false,
		"int":    0,
		"string": "",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.False(in.False)
	ast.Equal(0, in.Int)
	ast.Equal("", in.String)
}

func TestUnmarshalUsingDifferentKeys(t *testing.T) {
	type inner struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `bson:"string"`
	}
	m := map[string]interface{}{
		"no":     false,
		"int":    9,
		"string": "value",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.False(in.False)
	ast.Equal(9, in.Int)
	ast.True(len(in.String) == 0)
}

func TestUnmarshalNumberRangeInt(t *testing.T) {
	type inner struct {
		Value1  int    `key:"value1,range=[1:]"`
		Value2  int8   `key:"value2,range=[1:5]"`
		Value3  int16  `key:"value3,range=[1:5]"`
		Value4  int32  `key:"value4,range=[1:5]"`
		Value5  int64  `key:"value5,range=[1:5]"`
		Value6  uint   `key:"value6,range=[:5]"`
		Value8  uint8  `key:"value8,range=[1:5],string"`
		Value9  uint16 `key:"value9,range=[1:5],string"`
		Value10 uint32 `key:"value10,range=[1:5],string"`
		Value11 uint64 `key:"value11,range=[1:5],string"`
	}
	m := map[string]interface{}{
		"value1":  10,
		"value2":  int8(1),
		"value3":  int16(2),
		"value4":  int32(4),
		"value5":  int64(5),
		"value6":  uint(0),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(10, in.Value1)
	ast.Equal(int8(1), in.Value2)
	ast.Equal(int16(2), in.Value3)
	ast.Equal(int32(4), in.Value4)
	ast.Equal(int64(5), in.Value5)
	ast.Equal(uint(0), in.Value6)
	ast.Equal(uint8(1), in.Value8)
	ast.Equal(uint16(2), in.Value9)
	ast.Equal(uint32(4), in.Value10)
	ast.Equal(uint64(5), in.Value11)
}

func TestUnmarshalNumberRangeJsonNumber(t *testing.T) {
	type inner struct {
		Value3 uint   `key:"value3,range=(1:5]"`
		Value4 uint8  `key:"value4,range=(1:5]"`
		Value5 uint16 `key:"value5,range=(1:5]"`
	}
	m := map[string]interface{}{
		"value3": json.Number("2"),
		"value4": json.Number("4"),
		"value5": json.Number("5"),
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(uint(2), in.Value3)
	ast.Equal(uint8(4), in.Value4)
	ast.Equal(uint16(5), in.Value5)

	type inner1 struct {
		Value int `key:"value,range=(1:5]"`
	}
	m = map[string]interface{}{
		"value": json.Number("a"),
	}

	var in1 inner1
	ast.NotNil(UnmarshalKey(m, &in1))
}

func TestUnmarshalNumberRangeIntLeftExclude(t *testing.T) {
	type inner struct {
		Value3  uint   `key:"value3,range=(1:5]"`
		Value4  uint32 `key:"value4,default=4,range=(1:5]"`
		Value5  uint64 `key:"value5,range=(1:5]"`
		Value9  int    `key:"value9,range=(1:5],string"`
		Value10 int    `key:"value10,range=(1:5],string"`
		Value11 int    `key:"value11,range=(1:5],string"`
	}
	m := map[string]interface{}{
		"value3":  uint(2),
		"value4":  uint32(4),
		"value5":  uint64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(uint(2), in.Value3)
	ast.Equal(uint32(4), in.Value4)
	ast.Equal(uint64(5), in.Value5)
	ast.Equal(2, in.Value9)
	ast.Equal(4, in.Value10)
	ast.Equal(5, in.Value11)
}

func TestUnmarshalNumberRangeIntRightExclude(t *testing.T) {
	type inner struct {
		Value2  uint   `key:"value2,range=[1:5)"`
		Value3  uint8  `key:"value3,range=[1:5)"`
		Value4  uint16 `key:"value4,range=[1:5)"`
		Value8  int    `key:"value8,range=[1:5),string"`
		Value9  int    `key:"value9,range=[1:5),string"`
		Value10 int    `key:"value10,range=[1:5),string"`
	}
	m := map[string]interface{}{
		"value2":  uint(1),
		"value3":  uint8(2),
		"value4":  uint16(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(uint(1), in.Value2)
	ast.Equal(uint8(2), in.Value3)
	ast.Equal(uint16(4), in.Value4)
	ast.Equal(1, in.Value8)
	ast.Equal(2, in.Value9)
	ast.Equal(4, in.Value10)
}

func TestUnmarshalNumberRangeIntExclude(t *testing.T) {
	type inner struct {
		Value3  int `key:"value3,range=(1:5)"`
		Value4  int `key:"value4,range=(1:5)"`
		Value9  int `key:"value9,range=(1:5),string"`
		Value10 int `key:"value10,range=(1:5),string"`
	}
	m := map[string]interface{}{
		"value3":  2,
		"value4":  4,
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(2, in.Value3)
	ast.Equal(4, in.Value4)
	ast.Equal(2, in.Value9)
	ast.Equal(4, in.Value10)
}

func TestUnmarshalNumberRangeIntOutOfRange(t *testing.T) {
	type inner1 struct {
		Value int64 `key:"value,default=3,range=(1:5)"`
	}

	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": json.Number("6"),
	}, &in1))

	type inner2 struct {
		Value int64 `key:"value,optional,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(5),
	}, &in2))

	type inner3 struct {
		Value int64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(6),
	}, &in3))

	type inner4 struct {
		Value int64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(6),
	}, &in4))
}

func TestUnmarshalNumberRangeFloat(t *testing.T) {
	type inner struct {
		Value2  float32 `key:"value2,range=[1:5]"`
		Value3  float32 `key:"value3,range=[1:5]"`
		Value4  float64 `key:"value4,range=[1:5]"`
		Value5  float64 `key:"value5,range=[1:5]"`
		Value8  float64 `key:"value8,range=[1:5],string"`
		Value9  float64 `key:"value9,range=[1:5],string"`
		Value10 float64 `key:"value10,range=[1:5],string"`
		Value11 float64 `key:"value11,range=[1:5],string"`
	}
	m := map[string]interface{}{
		"value2":  float32(1),
		"value3":  float32(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(float32(1), in.Value2)
	ast.Equal(float32(2), in.Value3)
	ast.Equal(float64(4), in.Value4)
	ast.Equal(float64(5), in.Value5)
	ast.Equal(float64(1), in.Value8)
	ast.Equal(float64(2), in.Value9)
	ast.Equal(float64(4), in.Value10)
	ast.Equal(float64(5), in.Value11)
}

func TestUnmarshalNumberRangeFloatLeftExclude(t *testing.T) {
	type inner struct {
		Value3  float64 `key:"value3,range=(1:5]"`
		Value4  float64 `key:"value4,range=(1:5]"`
		Value5  float64 `key:"value5,range=(1:5]"`
		Value9  float64 `key:"value9,range=(1:5],string"`
		Value10 float64 `key:"value10,range=(1:5],string"`
		Value11 float64 `key:"value11,range=(1:5],string"`
	}
	m := map[string]interface{}{
		"value3":  float64(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(float64(2), in.Value3)
	ast.Equal(float64(4), in.Value4)
	ast.Equal(float64(5), in.Value5)
	ast.Equal(float64(2), in.Value9)
	ast.Equal(float64(4), in.Value10)
	ast.Equal(float64(5), in.Value11)
}

func TestUnmarshalNumberRangeFloatRightExclude(t *testing.T) {
	type inner struct {
		Value2  float64 `key:"value2,range=[1:5)"`
		Value3  float64 `key:"value3,range=[1:5)"`
		Value4  float64 `key:"value4,range=[1:5)"`
		Value8  float64 `key:"value8,range=[1:5),string"`
		Value9  float64 `key:"value9,range=[1:5),string"`
		Value10 float64 `key:"value10,range=[1:5),string"`
	}
	m := map[string]interface{}{
		"value2":  float64(1),
		"value3":  float64(2),
		"value4":  float64(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(float64(1), in.Value2)
	ast.Equal(float64(2), in.Value3)
	ast.Equal(float64(4), in.Value4)
	ast.Equal(float64(1), in.Value8)
	ast.Equal(float64(2), in.Value9)
	ast.Equal(float64(4), in.Value10)
}

func TestUnmarshalNumberRangeFloatExclude(t *testing.T) {
	type inner struct {
		Value3  float64 `key:"value3,range=(1:5)"`
		Value4  float64 `key:"value4,range=(1:5)"`
		Value9  float64 `key:"value9,range=(1:5),string"`
		Value10 float64 `key:"value10,range=(1:5),string"`
	}
	m := map[string]interface{}{
		"value3":  float64(2),
		"value4":  float64(4),
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalKey(m, &in))
	ast.Equal(float64(2), in.Value3)
	ast.Equal(float64(4), in.Value4)
	ast.Equal(float64(2), in.Value9)
	ast.Equal(float64(4), in.Value10)
}

func TestUnmarshalNumberRangeFloatOutOfRange(t *testing.T) {
	type inner1 struct {
		Value float64 `key:"value,range=(1:5)"`
	}

	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": json.Number("6"),
	}, &in1))

	type inner2 struct {
		Value float64 `key:"value,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(5),
	}, &in2))

	type inner3 struct {
		Value float64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(6),
	}, &in3))

	type inner4 struct {
		Value float64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(6),
	}, &in4))
}

func TestUnmarshalRangeError(t *testing.T) {
	type inner1 struct {
		Value int `key:",range="`
	}
	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in1))

	type inner2 struct {
		Value int `key:",range=["`
	}
	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in2))

	type inner3 struct {
		Value int `key:",range=[:"`
	}
	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in3))

	type inner4 struct {
		Value int `key:",range=[:]"`
	}
	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in4))

	type inner5 struct {
		Value int `key:",range={:]"`
	}
	var in5 inner5
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in5))

	type inner6 struct {
		Value int `key:",range=[:}"`
	}
	var in6 inner6
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in6))

	type inner7 struct {
		Value int `key:",range=[]"`
	}
	var in7 inner7
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in7))

	type inner8 struct {
		Value int `key:",range=[a:]"`
	}
	var in8 inner8
	assert.NotNil(t, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in8))

	type inner9 struct {
		Value int `key:",range=[:a]"`
	}
	var in9 inner9
	assert.NotNil(t, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in9))

	type inner10 struct {
		Value int `key:",range"`
	}
	var in10 inner10
	assert.NotNil(t, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in10))

	type inner11 struct {
		Value int `key:",range=[1,2]"`
	}
	var in11 inner11
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": "a",
	}, &in11))
}

func TestUnmarshalNestedMap(t *testing.T) {
	var c struct {
		Anything map[string]map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]map[string]interface{}{
			"inner": {
				"id":   "1",
				"name": "any",
			},
		},
	}

	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &c))
	assert.Equal(t, "1", c.Anything["inner"]["id"])
}

func TestUnmarshalNestedMapMismatch(t *testing.T) {
	var c struct {
		Anything map[string]map[string]map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]map[string]interface{}{
			"inner": {
				"name": "any",
			},
		},
	}

	assert.NotNil(t, NewUnmarshaler("json").Unmarshal(m, &c))
}

func TestUnmarshalNestedMapSimple(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]interface{}{
			"id":   "1",
			"name": "any",
		},
	}

	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &c))
	assert.Equal(t, "1", c.Anything["id"])
}

func TestUnmarshalNestedMapSimpleTypeMatch(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]string{
			"id":   "1",
			"name": "any",
		},
	}

	assert.Nil(t, NewUnmarshaler("json").Unmarshal(m, &c))
	assert.Equal(t, "1", c.Anything["id"])
}

func TestUnmarshalValuer(t *testing.T) {
	unmarshaler := NewUnmarshaler(jsonTagKey)
	var foo string
	err := unmarshaler.UnmarshalValuer(nil, foo)
	assert.NotNil(t, err)
}

func BenchmarkUnmarshalString(b *testing.B) {
	type inner struct {
		Value string `key:"value"`
	}
	m := map[string]interface{}{
		"value": "first",
	}

	for i := 0; i < b.N; i++ {
		var in inner
		if err := UnmarshalKey(m, &in); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalStruct(b *testing.B) {
	b.ReportAllocs()

	m := map[string]interface{}{
		"Ids": []map[string]interface{}{
			{
				"First":  1,
				"Second": 2,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		var v struct {
			Ids []struct {
				First  int
				Second int
			}
		}
		if err := UnmarshalKey(m, &v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapToStruct(b *testing.B) {
	data := map[string]interface{}{
		"valid": "1",
		"age":   "5",
		"name":  "liao",
	}
	type anonymous struct {
		Valid bool
		Age   int
		Name  string
	}

	for i := 0; i < b.N; i++ {
		var an anonymous
		if valid, ok := data["valid"]; ok {
			an.Valid = valid == "1"
		}
		if age, ok := data["age"]; ok {
			ages, _ := age.(string)
			an.Age, _ = strconv.Atoi(ages)
		}
		if name, ok := data["name"]; ok {
			names, _ := name.(string)
			an.Name = names
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	data := map[string]interface{}{
		"valid": "1",
		"age":   "5",
		"name":  "liao",
	}
	type anonymous struct {
		Valid bool   `key:"valid,string"`
		Age   int    `key:"age,string"`
		Name  string `key:"name"`
	}

	for i := 0; i < b.N; i++ {
		var an anonymous
		UnmarshalKey(data, &an)
	}
}
