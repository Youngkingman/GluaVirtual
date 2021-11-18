package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
	"github.com/Youngkingman/GluaVirtual/luaState/state"
	vm "github.com/Youngkingman/GluaVirtual/virtualMachine"
)

func main() {
	data, err := ioutil.ReadFile("fornum.out")
	if err != nil {
		panic(err)
	}
	proto := binarychunk.Undump(data)
	LuaEntry(proto)
}

func LuaEntry(proto *binarychunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	st := state.New(nRegs+8, proto)
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
