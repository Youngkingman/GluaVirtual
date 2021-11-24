package state

import (
	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
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
		if c.proto != nil {
			//打印调试信息
			// fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
			// 	c.proto.LineDefined, c.proto.LastLineDefined)
			//实际执行Lua函数
			st.callLuaClosure(nArgs, nResults, c)
		} else {
			//实际执行Go函数
			st.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not a function to call")
	}
}

func (st *LuaState) callLuaClosure(nArgs, nResults int, c *luaClosure) {
	//根据寄存器数量创建新的luaStack（栈帧）
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStk := newLuaStack(nRegs+luaApi.LUA_MINSTACK, st) //稍微多一点点栈空间
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

func (st *LuaState) callGoClosure(nArgs, nResults int, c *luaClosure) {
	//根据所需寄存器数量创建luaStack(栈帧)
	newStk := newLuaStack(nArgs+luaApi.LUA_MINSTACK, st)
	newStk.closure = c

	//参数加入新栈帧
	args := st.stack.popN(nArgs)
	newStk.pushN(args, nArgs)
	st.stack.pop() //本栈帧丢弃go函数

	//换栈帧，执行go函数,换回来
	st.pushLuaStack(newStk)
	r := c.goFunc(st)
	st.popLuaStack()

	//如果需要返回值，那就返回
	if nResults != 0 {
		results := newStk.popN(r)
		st.stack.check(len(results))
		st.stack.pushN(results, nResults)
	}
}
