package state

func (st *LuaState) PC() int {
	return st.pc
}

func (st *LuaState) OffsetPC(n int) {
	st.pc += n
}

func (st *LuaState) Fetch() uint32 {
	i := st.proto.Code[st.pc]
	st.pc++
	return i
}

func (st *LuaState) GetConst(idx int) {
	c := st.proto.Constants[idx]
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
