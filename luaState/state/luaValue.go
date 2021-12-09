package state

import (
	"fmt"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
	"github.com/Youngkingman/GluaVirtual/numTrans"
)

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *luaClosure:
		return LUA_TFUNCTION
	default:
		panic("fuck type")
	}
}

func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func converToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return numTrans.ParseFloat(x)
	default:
	}
	return 0, false
}

func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return numTrans.Float2Integer(x)
	case string:
		return numTrans.String2Integer(x)
	default:
	}
	return 0, false
}

func setMetatable(val luaValue, mt *luaTable, st *LuaState) {
	if t, ok := val.(*luaTable); ok {
		//赋予元表,直接修改元表字段
		t.metatable = mt
		return
	}
	//不是表类型，根据类型放在注册表里,元表若为nil效果相当于删除元表
	key := fmt.Sprintf("_MT%d", typeOf(val))
	st.registry.put(key, mt)
}

func getMetatable(val luaValue, st *LuaState) *luaTable {
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := st.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}
	return nil
}

func callMetamethod(a, b luaValue, mmName string, st *LuaState) (luaValue, bool) {
	var mm luaValue
	if mm = getMetafield(b, mmName, st); mm == nil {
		return nil, false
	}

	st.stack.check(4)
	st.stack.push(mm)
	st.stack.push(a)
	st.stack.push(b)
	st.Call(2, 1)
	return st.stack.pop(), true
}

func getMetafield(val luaValue, fieldName string, st *LuaState) luaValue {
	if mt := getMetatable(val, st); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}
