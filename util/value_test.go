package luautil

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestCleanStack(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	l := lua.NewState()
	require.NotNil(l)
	defer l.Close()

	err := l.DoString(`
		return 1, 2
	`)
	require.NoError(err)
	CleanStack(l)
	assert.Equal(0, l.GetTop())

	err = l.DoString(`
		local a = 1
	`)
	require.NoError(err)
	CleanStack(l)
	assert.Equal(0, l.GetTop())
}

func TestGetArgs(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	l := lua.NewState()
	require.NotNil(l)
	defer l.Close()

	err := l.DoString(`
		return 1, 2
	`)
	require.NoError(err)
	v := GetArgs(l, 1)
	assert.Equal([]interface{}{1, 2}, v)
}

func TestGetValue(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	l := lua.NewState()
	require.NotNil(l)
	defer l.Close()

	b := bytes.NewBufferString(`
		return {
			a = 1,
			b = 1.1,
			c = "foo",
			d = {
				e = "baz",
			},
			f = {"buz", 5},
			g = {},
			h = true,
		}
	`)
	fn, err := l.Load(b, "")
	require.NoError(err)
	l.Push(fn)
	l.Call(0, 1)
	i := GetValue(l, l.Get(-1))
	assert.Equal(map[string]interface{}{
		"a": 1,
		"b": 1.1,
		"c": "foo",
		"d": map[string]interface{}{
			"e": "baz",
		},
		"f": []interface{}{"buz", 5},
		"g": []interface{}{},
		"h": true,
	}, i)
}

func TestToLuaValue(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	l := lua.NewState()
	require.NotNil(l)
	defer l.Close()

	i := int64(1)
	lv := ToLuaValue(l, i)
	assert.Equal(lua.LTNumber, lv.Type())
	assert.Equal("1", lv.String())

	j := true
	lv = ToLuaValue(l, j)
	assert.Equal(lua.LTBool, lv.Type())
	assert.Equal("true", lv.String())

	lv = ToLuaValue(l, nil)
	assert.Equal(lua.LTNil, lv.Type())
	assert.Equal("nil", lv.String())
}

func checkToLuaTable(t *testing.T, f func(*lua.LState, reflect.Value) lua.LValue, i interface{}, code string) {
	l := lua.NewState()
	require.NotNil(t, l)
	defer l.Close()
	initialStackSize := l.GetTop()

	v := f(l, reflect.ValueOf(i))
	l.SetGlobal("ctx", v)
	assert.Equal(t, initialStackSize, l.GetTop())

	b := bytes.NewBufferString(code)
	fn, err := l.Load(b, "")
	require.Nil(t, err)
	l.Push(fn)
	l.Call(0, 1)
	assert.True(t, l.ToBool(-1))
	l.Remove(-1)
	assert.Equal(t, initialStackSize, l.GetTop())
}

func TestTableFromStruct(t *testing.T) {
	type Foo struct {
		A int
		B string
	}

	type Bar struct {
		C Foo
		D bool `lua:"d"`
	}

	type Baz struct {
		Bar `lua:",inline"`
		E   *string
		F   int `lua:"-"`
		G   struct{}
		H   *int
		i   int
	}

	str := "wut"
	i := Baz{Bar{Foo{1, "wat"}, true}, &str, 5, struct{}{}, nil, 6}
	checkToLuaTable(t, luaTableFromStruct, i, `
		if ctx.C.A ~= 1 then return false end
		if ctx.C.B ~= "wat" then return false end
		if ctx.d ~= true then return false end
		if ctx.E ~= "wut" then return false end
		if ctx.F ~= nil then return false end
		if next(ctx.G) ~= nil then return false end
		if ctx.H ~= nil then return false end
		if ctx.i ~= nil then return false end
		return true
	`)

	type Qux struct {
		Foo `lua:"test"`
	}

	j := Qux{Foo{1, "2"}}
	checkToLuaTable(t, luaTableFromStruct, j, `
		if ctx.test.A ~= 1 then return false end
		if ctx.test.B ~= "2" then return false end
		return true
	`)
}

func TestTableFromMap(t *testing.T) {
	m := map[interface{}]interface{}{
		"A": 1,
		5:   "FOO",
		true: map[string]interface{}{
			"foo": "bar",
		},
	}
	checkToLuaTable(t, luaTableFromMap, m, `
		if ctx.A ~= 1 then return false end
		if ctx[5] ~= "FOO" then return false end
		if ctx[true].foo ~= "bar" then return false end
		return true
	`)
}

func TestTableFromSlice(t *testing.T) {
	s := []interface{}{
		"foo",
		true,
		4,
		[]string{
			"bar",
			"baz",
		},
		[2]int{1, 2},
	}
	checkToLuaTable(t, luaTableFromSlice, s, `
		if ctx[1] ~= "foo" then return false end
		if ctx[2] ~= true then return false end
		if ctx[3] ~= 4 then return false end
		if ctx[4][1] ~= "bar" then return false end
		if ctx[4][2] ~= "baz" then return false end
		if ctx[5][1] ~= 1 then return false end
		if ctx[5][2] ~= 2 then return false end
		return true
	`)
}
