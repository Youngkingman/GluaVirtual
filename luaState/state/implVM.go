package state

func (st *LuaState) PC() int {
	return st.stack.pc
}

func (st *LuaState) OffsetPC(n int) {
	st.stack.pc += n
}

func (st *LuaState) Fetch() uint32 {
	i := st.stack.closure.proto.Code[st.stack.pc]
	st.stack.pc++
	return i
}

func (st *LuaState) GetConst(idx int) {
	c := st.stack.closure.proto.Constants[idx]
	st.stack.push(c)
}

//rk is tpye OpArgK
func (st *LuaState) GetRK(rk int) {
	if rk > 0xFF { //a constant index
		st.GetConst(rk & 0xFF)
	} else { //a register
		st.PushValue(rk + 1)
	}
}
