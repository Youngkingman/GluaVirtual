package state

func (st *luaState) PushNil() {
	st.stack.push(nil)
}

func (st *luaState) PushBoolean(b bool) {
	st.stack.push(b)
}

func (st *luaState) PushInteger(n int64) {
	st.stack.push(n)
}

func (st *luaState) PushNumber(n float64) {
	st.stack.push(n)
}

func (st *luaState) PushString(s string) {
	st.stack.push(s)
}
