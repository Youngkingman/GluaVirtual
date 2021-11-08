package state

import (
	"fmt"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

//series of type information method
func (st *luaState) TypeName(tp LuaType) string {
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

func (st *luaState) Type(idx int) LuaType {
	if st.stack.isValid(idx) {
		val := st.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (st *luaState) IsNone(idx int) bool {
	return st.Type(idx) == LUA_TNONE
}

func (st *luaState) IsNil(idx int) bool {
	return st.Type(idx) == LUA_TNIL
}

func (st *luaState) IsNoneOrNil(idx int) bool {
	return st.Type(idx) <= LUA_TNIL
}

func (st *luaState) IsBoolean(idx int) bool {
	return st.Type(idx) == LUA_TBOOLEAN
}

func (st *luaState) IsTable(idx int) bool {
	return st.Type(idx) == LUA_TTABLE
}

func (st *luaState) IsFunction(idx int) bool {
	return st.Type(idx) == LUA_TFUNCTION
}

func (st *luaState) IsThread(idx int) bool {
	return st.Type(idx) == LUA_TTHREAD
}

func (st *luaState) IsString(idx int) bool {
	t := st.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (st *luaState) IsNumber(idx int) bool {
	_, ok := st.ToNumberX(idx)
	return ok
}

func (st *luaState) IsInteger(idx int) bool {
	val := st.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (st *luaState) ToBoolean(idx int) bool {
	val := st.stack.get(idx)
	return convertToBoolean(val)
}

func (st *luaState) ToInteger(idx int) int64 {
	i, _ := st.ToIntegerX(idx)
	return i
}

func (st *luaState) ToIntegerX(idx int) (int64, bool) {
	val := st.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

func (st *luaState) ToNumber(idx int) float64 {
	n, _ := st.ToNumberX(idx)
	return n
}

func (st *luaState) ToNumberX(idx int) (float64, bool) {
	val := st.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (st *luaState) ToString(idx int) string {
	s, _ := st.ToStringX(idx)
	return s
}

func (st *luaState) ToStringX(idx int) (string, bool) {
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
