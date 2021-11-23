package state

import (
	"fmt"

	"github.com/Youngkingman/GluaVirtual/binarychunk"
	vm "github.com/Youngkingman/GluaVirtual/virtualMachine"
)

func (st *LuaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binarychunk.Undump(chunk)
	c := newLuaClosure(proto)
	st.stack.push(c)
	return 0
}

func (st *LuaState) Call(nArgs, nResults int) {
	val := st.stack.get(-(nArgs + 1))
	if c, ok := val.(*luaClosure); ok {
		//打印调试信息
		fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		//实际执行函数
		st.callLuaClosure(nArgs, nResults, c)
	} else {
		panic("not a function to call")
	}
}

func (st *LuaState) callLuaClosure(nArgs, nResults int, c *luaClosure) {
	//根据寄存器数量创建新的luaStack
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStk := newLuaStack(nRegs + 20) //稍微多一点点栈空间
	newStk.closure = c

	//旧luaStack的调用参数传出，第一个参数为函数，其余为函数参数
	funcArgs := st.stack.popN(nArgs + 1)
	newStk.pushN(funcArgs[1:], nParams)
	newStk.top = nRegs
	//记录可变变量
	if nArgs > nParams && isVararg {
		newStk.varargs = funcArgs[nParams+1:]
	}

	//换栈,递归执行，退出
	st.pushLuaStack(newStk)
	st.runluaClosure()
	st.popLuaStack()

	if nResults != 0 {
		results := newStk.popN(newStk.top - nRegs)
		st.stack.check(len(results)) //扩展空间
	}
}

//事实上执行当前状态机的闭包
func (st *LuaState) runluaClosure() {
	for {
		inst := vm.Instruction(st.Fetch())
		inst.Execute(st)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}
