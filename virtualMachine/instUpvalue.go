package vm

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

func getTabUp(i Instruction, vm luaApi.LuaVMInterface) {
	a, _, c := i.ABC()
	a += 1

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
	vm.Pop(1)
}
