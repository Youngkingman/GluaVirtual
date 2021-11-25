package vm

import (
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

//对数值型的for循环，lua编译器会预先使用三个特殊局部变量，分别存放数值、循环限制条件和循环步长
//对应寄存器为指令中的a/a+1/a+2，自定义的循环变量i则存放在a+3里面

//ForPrep,R(A) += R(A+2);PC += sbx
func forPrep(i Instruction, vm luaApi.LuaVMInterface) {
	a, sBx := i.AsBx()
	a += 1

	if vm.Type(a) == luaApi.LUA_TSTRING {
		//将字符串解析为数值型，下同
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}
	if vm.Type(a+1) == luaApi.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}
	if vm.Type(a+2) == luaApi.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}

	//正式循环前减去步长
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(luaApi.LUA_OPSUB)
	vm.Replace(a)
	//跳转到forLoop执行循环
	vm.OffsetPC(sBx)
}

//ForLoop, R(A) += R(A+2)
//if R(A) <?=R(A+1) then {pc += sbx;R(A+3) = R(A)}
func forLoop(i Instruction, vm luaApi.LuaVMInterface) {
	a, sBx := i.AsBx()
	a += 1
	vm.PushValue(a + 2)
	vm.PushValue(a)
	//先加上步长
	vm.Arith(luaApi.LUA_OPADD)
	vm.Replace(a)

	isPoitiveStep := vm.ToNumber(a+2) >= 0
	if isPoitiveStep && vm.Compare(a, a+1, luaApi.LUA_OPLE) || !isPoitiveStep && vm.Compare(a+1, a, luaApi.LUA_OPLE) {
		//条件成立，执行循环
		vm.OffsetPC(sBx)
		//更新循环变量
		vm.Copy(a, a+3)
	}
}
