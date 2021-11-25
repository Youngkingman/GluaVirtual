package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
	"github.com/Youngkingman/GluaVirtual/luaState/state"
	vm "github.com/Youngkingman/GluaVirtual/virtualMachine"
)

var filenames = [...]string{
	"test.out",
	"fornum.out",
	"funcCall.out",
	"hw.out",
}

func Test_ParseFunc(t *testing.T) {
	data, err := ioutil.ReadFile(filenames[0])
	if err != nil {
		panic(err)
	}
	proto := binarychunk.Undump(data)
	list(proto)
}

func Test_ExcuteOpt(t *testing.T) {
	data, err := ioutil.ReadFile(filenames[1])
	if err != nil {
		panic(err)
	}
	proto := binarychunk.Undump(data)
	LuaEntry(proto)
}

func Test_FunctionCall(t *testing.T) {
	data, err := ioutil.ReadFile(filenames[1])
	if err != nil {
		panic(err)
	}
	st := state.New()
	st.Load(data, filenames[1], "b")
	st.Call(0, 0)
}

func Test_Print(t *testing.T) {
	data, err := ioutil.ReadFile(filenames[3])
	if err != nil {
		panic(err)
	}
	st := state.New()
	st.Register("print", print)
	st.Load(data, "chunk", "b")
	st.Call(0, 0)
}

func print(st luaApi.LuaStateInterface) int {
	nArgs := st.GetTop()
	for i := 1; i <= nArgs; i++ {
		if st.IsBoolean(i) {
			fmt.Printf("%t", st.ToBoolean(i))
		} else {
			fmt.Print(st.TypeName(st.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
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

func printStack(ls *state.LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case luaApi.LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case luaApi.LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case luaApi.LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func LuaEntry(proto *binarychunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	st := state.New()
	st.SetTop(nRegs)
	for {
		pc := st.PC()
		inst := vm.Instruction(st.Fetch())
		if inst.Opcode() != vm.OP_RETURN {
			inst.Execute(st)

			fmt.Printf("[%02d] %s", pc+1, inst.OpName())
			printStack(st)
		} else {
			break
		}
	}
}
