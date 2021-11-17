package vm

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

//栈索引需要在源索引基础上加一
func move(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a)
}

func jmp(i Instruction, vm luaApi.LuaVMInterface) {
	a, sbx := i.AsBx()
	vm.OffsetPC(sbx)
	if a != 0 {
		//和upvalue有关
		panic("")
	}
}
