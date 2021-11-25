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

//返回堆栈大小即寄存器数目
func (st *LuaState) RegisterCount() int {
	return int(st.stack.closure.proto.MaxStackSize)
}

//变长参数加入堆栈，多退少补
func (st *LuaState) LoadVararg(n int) {
	if n < 0 {
		n = len(st.stack.varargs)
	}
	st.stack.check(n)
	st.stack.pushN(st.stack.varargs, n)
}

//加载函数闭包再推入堆栈
func (st *LuaState) LoadProto(idx int) {
	proto := st.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	st.stack.push(closure)
}
