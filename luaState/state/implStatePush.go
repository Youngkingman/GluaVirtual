package state

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
