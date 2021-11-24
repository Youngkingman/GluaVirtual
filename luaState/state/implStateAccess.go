package state

import (
	"fmt"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

//series of type information method
func (st *LuaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (st *LuaState) Type(idx int) LuaType {
	if st.stack.isValid(idx) {
		val := st.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (st *LuaState) IsNone(idx int) bool {
	return st.Type(idx) == LUA_TNONE
}

func (st *LuaState) IsNil(idx int) bool {
	return st.Type(idx) == LUA_TNIL
}

func (st *LuaState) IsNoneOrNil(idx int) bool {
	return st.Type(idx) <= LUA_TNIL
}

func (st *LuaState) IsBoolean(idx int) bool {
	return st.Type(idx) == LUA_TBOOLEAN
}

func (st *LuaState) IsTable(idx int) bool {
	return st.Type(idx) == LUA_TTABLE
}

func (st *LuaState) IsFunction(idx int) bool {
	return st.Type(idx) == LUA_TFUNCTION
}

func (st *LuaState) IsThread(idx int) bool {
	return st.Type(idx) == LUA_TTHREAD
}

func (st *LuaState) IsString(idx int) bool {
	t := st.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (st *LuaState) IsNumber(idx int) bool {
	_, ok := st.ToNumberX(idx)
	return ok
}

func (st *LuaState) IsInteger(idx int) bool {
	val := st.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (st *LuaState) ToBoolean(idx int) bool {
	val := st.stack.get(idx)
	return convertToBoolean(val)
}

func (st *LuaState) ToInteger(idx int) int64 {
	i, _ := st.ToIntegerX(idx)
	return i
}

func (st *LuaState) ToIntegerX(idx int) (int64, bool) {
	val := st.stack.get(idx)
	return convertToInteger(val)
}

func (st *LuaState) ToNumber(idx int) float64 {
	n, _ := st.ToNumberX(idx)
	return n
}

func (st *LuaState) ToNumberX(idx int) (float64, bool) {
	val := st.stack.get(idx)
	return converToFloat(val)
}

func (st *LuaState) ToString(idx int) string {
	s, _ := st.ToStringX(idx)
	return s
}

func (st *LuaState) ToStringX(idx int) (string, bool) {
	val := st.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x) // todo
		st.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (st *LuaState) IsGoFunction(idx int) bool {
	val := st.stack.get(idx)
	if c, ok := val.(*luaClosure); ok {
		return c.goFunc == nil
	}
	return false
}

func (st *LuaState) ToGoFunction(idx int) GoFunction {
	val := st.stack.get(idx)
	if c, ok := val.(*luaClosure); ok {
		return c.goFunc
	}
	return nil
}
