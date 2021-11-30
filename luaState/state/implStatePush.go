package state

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

func (st *LuaState) PushNil() {
	st.stack.push(nil)
}

func (st *LuaState) PushBoolean(b bool) {
	st.stack.push(b)
}

func (st *LuaState) PushInteger(n int64) {
	st.stack.push(n)
}

func (st *LuaState) PushNumber(n float64) {
	st.stack.push(n)
}

func (st *LuaState) PushString(s string) {
	st.stack.push(s)
}

func (st *LuaState) PushGoFunction(f luaApi.GoFunction) {
	st.stack.push(newGoClosure(f, 0))
}

func (st *LuaState) PushGlobalTable() {
	global := st.registry.get(luaApi.LUA_RIDX_GLOBALS)
	st.stack.push(global)
}
