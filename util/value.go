package luautil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	lua "github.com/yuin/gopher-lua"
)

// CleanStack from glua vm top
func CleanStack(l *lua.LState) {
	for l.GetTop() > 0 {
		l.Remove(-1)
	}
}

// GetArgs from glua vm top
func GetArgs(l *lua.LState, start int) []interface{} {
	var v []interface{}
	for i := start; i <= l.GetTop(); i++ {
		arg := GetValue(l, l.CheckAny(i))
		v = append(v, arg)
	}
	return v
}

// GetValue converts glua vm value to go value
func GetValue(l *lua.LState, v lua.LValue) interface{} {
	switch t := v.Type(); t {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return lua.LVAsBool(v)
	case lua.LTNumber:
		f := lua.LVAsNumber(v)
		if float64(f) == float64(int(f)) {
			return int(f)
		}
		return float64(f)
	case lua.LTString:
		return lua.LVAsString(v)
	case lua.LTTable:
		m := map[string]interface{}{}
		tb := v.(*lua.LTable)
		arrSize := 0
		tb.ForEach(func(k, val lua.LValue) {
			key := GetValue(l, k)
			if keyi, ok := key.(int); ok {
				if arrSize >= 0 && arrSize < keyi {
					arrSize = keyi
				}
				key = strconv.Itoa(keyi)
			} else {
				arrSize = -1
			}
			m[key.(string)] = GetValue(l, val)
		})

		if arrSize >= 0 {
			ms := make([]interface{}, arrSize)
			for i := 0; i < arrSize; i++ {
				ms[i] = m[strconv.Itoa(i+1)]
			}
			return ms
		}

		return m
	default:
		panic(fmt.Sprintf("unknown lua type: %s", t))
	}
}

// ToLuaValue converts go value to glua vm value
func ToLuaValue(l *lua.LState, i interface{}) lua.LValue {
	if i == nil {
		return lua.LNil
	}
	if v, ok := i.(lua.LValue); ok {
		return v
	}

	switch ii := i.(type) {
	case bool:
		return lua.LBool(ii)
	case int:
		return lua.LNumber(ii)
	case int8:
		return lua.LNumber(ii)
	case int16:
		return lua.LNumber(ii)
	case int32:
		return lua.LNumber(ii)
	case int64:
		return lua.LNumber(ii)
	case uint:
		return lua.LNumber(ii)
	case uint8:
		return lua.LNumber(ii)
	case uint16:
		return lua.LNumber(ii)
	case uint32:
		return lua.LNumber(ii)
	case uint64:
		return lua.LNumber(ii)
	case float64:
		return lua.LNumber(ii)
	case float32:
		return lua.LNumber(ii)
	case string:
		return lua.LString(ii)
	case []byte:
		return lua.LString(ii)
	default:
		v := reflect.ValueOf(i)
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				return lua.LNil
			}
			return ToLuaValue(l, v.Elem().Interface())

		case reflect.Struct:
			return luaTableFromStruct(l, v)

		case reflect.Map:
			return luaTableFromMap(l, v)

		case reflect.Slice:
			return luaTableFromSlice(l, v)

		case reflect.Array:
			return luaTableFromSlice(l, v)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return lua.LNumber(v.Int())

		default:
			panic(fmt.Sprintf("unknown type being pushed onto lua stack: %T %+v", i, i))
		}
	}
}

func luaTableFromStruct(l *lua.LState, v reflect.Value) lua.LValue {
	tb := l.NewTable()
	return luaTableFromStructInner(l, tb, v)
}

func luaTableFromStructInner(l *lua.LState, tb *lua.LTable, v reflect.Value) lua.LValue {
	t := v.Type()
	for j := 0; j < v.NumField(); j++ {
		var inline bool
		name := t.Field(j).Name
		if unicode.IsLower(rune(name[0])) {
			continue
		}
		if tag := t.Field(j).Tag.Get("lua"); tag != "" {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] == "-" {
				continue
			} else if tagParts[0] != "" {
				name = tagParts[0]
			}
			if len(tagParts) > 1 && tagParts[1] == "inline" {
				inline = true
			}
		}
		if inline {
			luaTableFromStructInner(l, tb, v.Field(j))
		} else {
			tb.RawSetString(name, ToLuaValue(l, v.Field(j).Interface()))
		}
	}
	return tb
}

func luaTableFromMap(l *lua.LState, v reflect.Value) lua.LValue {
	tb := l.NewTable()
	for _, k := range v.MapKeys() {
		tb.RawSet(ToLuaValue(l, k.Interface()),
			ToLuaValue(l, v.MapIndex(k).Interface()))
	}
	return tb
}

func luaTableFromSlice(l *lua.LState, v reflect.Value) lua.LValue {
	tb := l.NewTable()
	for j := 0; j < v.Len(); j++ {
		tb.RawSetInt(j+1, // because lua is 1-indexed
			ToLuaValue(l, v.Index(j).Interface()))
	}
	return tb
}
