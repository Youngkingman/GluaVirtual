package state

import (
	"fmt"

	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
	vm "github.com/Youngkingman/GluaVirtual/virtualMachine"
)

func (st *LuaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binarychunk.Undump(chunk)
	//加载主函数
	c := newLuaClosure(proto)
	st.stack.push(c)
	if len(proto.Upvalues) > 0 { //设置环境变量
		env := st.registry.get(luaApi.LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

func (st *LuaState) Call(nArgs, nResults int) {
	val := st.stack.get(-(nArgs + 1))
	c, ok := val.(*luaClosure)
	//如果参数不是closure，则查找其__call元方法
	if !ok {
		if mf := getMetafield(val, "__call", st); mf != nil {
			if c, ok = mf.(*luaClosure); ok {
				st.stack.push(val)
				st.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			//打印调试信息
			//实际执行Lua函数
			st.callLuaClosure(nArgs, nResults, c)
		} else {
			//实际执行Go函数
			st.callGoClosure(nArgs, nResults, c)
		}
	}
	panic(fmt.Sprintf("not a function to call args:%d results:%d and no metamethod __call", nArgs, nResults))
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
		st.stack.pushN(results, nResults)
	}
}

//事实上执行当前状态机的闭包
func (st *LuaState) runluaClosure() {
	for {
		inst := vm.Instruction(st.Fetch())
		fmt.Printf("%s", inst.OpName())
		printOperands(inst)
		printStack(st)
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

func (st *LuaState) Error() int {
	err := st.stack.pop()
	panic(err)
}

func (st *LuaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := st.stack
	status = luaApi.LUA_ERRRUN

	//catch error
	defer func() {
		if err := recover(); err != nil {
			for st.stack != caller {
				st.popLuaStack()
			}
			st.stack.push(err)
		}
	}()

	st.Call(nArgs, nResults)
	status = luaApi.LUA_OK
	return
}

/************************/
/*调试时打印堆栈和指令信息*/
/************************/
func printStack(st *LuaState) {
	top := st.GetTop()
	for i := 1; i <= top; i++ {
		t := st.Type(i)
		switch t {
		case luaApi.LUA_TBOOLEAN:
			fmt.Printf("[%t]", st.ToBoolean(i))
		case luaApi.LUA_TNUMBER:
			fmt.Printf("[%g]", st.ToNumber(i))
		case luaApi.LUA_TSTRING:
			fmt.Printf("[%q]", st.ToString(i))
		default: // other values
			fmt.Printf("[%s]", st.TypeName(t))
		}
	}
	println()
}

func list(f *binarychunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

func printHeader(f *binarychunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binarychunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := vm.Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printDetail(f *binarychunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *binarychunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

func printOperands(i vm.Instruction) {
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.ArgBMode() != vm.OpArgN { //operands is used
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF) //means constants index
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.ArgCMode() != vm.OpArgN { //operator is used
			if c > 0xff {
				fmt.Printf(" %d", -1-c&0xFF) //means constants index
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.ArgBMode() == vm.OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.ArgBMode() == vm.OpArgU {
			fmt.Printf(" %d", bx)
		}
	case vm.IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case vm.IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	}
}
