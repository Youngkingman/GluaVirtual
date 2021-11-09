package state

import (
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
