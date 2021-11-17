package vm

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

//将索引a+1开始的b个寄存器设为nil
//预执行阶段会把寄存器数量计算好保存在原型中，那么假定此时虚拟机已经保留了必要的栈空间而不考虑越界
func loadNil(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1

	vm.PushNil()
	for i := a; i < a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

func loadBool(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, c := i.ABC()
	a += 1
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.OffsetPC(1)
	}
}

func loadK(i Instruction, vm luaApi.LuaVMInterface) {
	a, bx := i.ABx()
	a += 1

	vm.GetConst(bx)
	vm.Replace(a)
}
