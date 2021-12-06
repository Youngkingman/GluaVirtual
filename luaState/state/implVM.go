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

//加载函数闭包再推入堆栈,主要用于指令closure的加载,即子函数调用
func (st *LuaState) LoadProto(idx int) {
	stk := st.stack
	subProto := stk.closure.proto.Protos[idx] //用于调用的子闭包
	closure := newLuaClosure(subProto)
	st.stack.push(closure)
	//根据proto中的Upvalues表来初始化闭包中的upvalue值
	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 { //upvalue属于当前函数,需要访问当前函数局部变量并映射closure
			if stk.openuvs == nil {
				stk.openuvs = map[int]*upvalue{}
			}
			if openuv, found := stk.openuvs[uvIdx]; found { //局部变量仍在栈上没有退出作用域
				closure.upvals[i] = openuv
			} else {
				closure.upvals[i] = &upvalue{&stk.slots[uvIdx]}
				stk.openuvs[uvIdx] = closure.upvals[i]
			}
		} else { //upvalue属于外围函数，直接将当前stack所属的closure的upvalue传递至下层closure
			closure.upvals[i] = stk.closure.upvals[uvIdx]
		}
	}
}

func (st *LuaState) CloseUpvalues(a int) {
	for i, openuv := range st.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(st.stack.openuvs, i)
		}
	}
}
