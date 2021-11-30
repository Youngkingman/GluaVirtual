package vm

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

//将指定的upvalue复制到寄存器（栈）上,栈伪索引从1开始，指令中索引从0开始
func getUpval(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(luaApi.LuaUpvalueIndex(b), a)
}

//getUpval的反过程
func setUpval(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, luaApi.LuaUpvalueIndex(b))
}

//upvalue如果是表,则不必整个复制只要取索引处值即可
//寄存器由A指定，upvalue表索引由操作数B指定，表Key由操作数C指定
func getTabUp(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(luaApi.LuaUpvalueIndex(b))
	vm.Replace(a)
}

//getTabUp的反过程
func setTabUp(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(luaApi.LuaUpvalueIndex(a))
}
