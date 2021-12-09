package state

import (
	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

func (st *LuaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := st.stack.get(idx1)
	b := st.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return _eq(a, b, st)
	case LUA_OPLT:
		return _lt(a, b, st)
	case LUA_OPLE:
		return _le(a, b, st)
	}
	return false
}

func (st *LuaState) RawEqual(idx1, idx2 int) bool {
	if !st.stack.isValid(idx1) || !st.stack.isValid(idx2) {
		return false
	}

	a := st.stack.get(idx1)
	b := st.stack.get(idx2)
	return _eq(a, b, nil)
}

func _eq(a, b luaValue, st *LuaState) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && st != nil { //为nil时不希望执行元方法
			if result, ok := callMetamethod(x, y, "__eq", st); ok {
				return convertToBoolean(result)
			}
		}
	default:
		return a == b
	}
	return true
}

func _lt(a, b luaValue, st *LuaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}
	if result, ok := callMetamethod(a, b, "__lt", st); ok {
		return convertToBoolean(result)
	} else {
		panic("comparison error!")
	}

}

func _le(a, b luaValue, st *LuaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	if result, ok := callMetamethod(a, b, "__le", st); ok {
		return convertToBoolean(result)
	} else if result, ok := callMetamethod(a, b, "__lt", st); ok { //没有le元方法就尝试下别的
		return !convertToBoolean(result)
	}
	panic("comparison error!")
}
