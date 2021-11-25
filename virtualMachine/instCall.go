package vm

import (
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

//R(A)=closure(KPROTO[Bx]) 实例化闭包为函数对象
func closure(i Instruction, vm luaApi.LuaVMInterface) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

//R(A),...,R(A+C-2) := R(A)(R(A+1),...,R(A+B-1)),c入参数+1，b返回值数+1
func call(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, c := i.ABC()
	a += 1

	nArgs := _pushFuncAndArgs(a, b, vm) //将被调函数和参数值推入堆栈
	vm.Call(nArgs, c-1)                 //进行函数本身的执行
	_popResults(a, c, vm)               //把返回值移动到适当的寄存器
}

//b为0表示要将参数全部展开传入，用于函数作为闭包参数时
func _pushFuncAndArgs(a, b int, vm luaApi.LuaVMInterface) (nArgs int) {
	if b >= 1 { //b-1 args
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else { //can be zero
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

//c为0表示要将参数全部展开传出，用于函数作为闭包返回时
func _popResults(a, c int, vm luaApi.LuaVMInterface) {
	if c == 1 { //have no result

	} else if c > 1 { //have c-1 result
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func _fixStack(a int, vm luaApi.LuaVMInterface) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

//return R(A),...,R(A+B-2)
func _return(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 { //no return values

	} else if b > 1 { //b-1 return values
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		_fixStack(a, vm)
	}
}

//R(A),R(A+1),...,R(A+B-2) = vararg
func vararg(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 {
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

//return R(A)(R(A+1),...,R(A+B-1))，用于尾递归优化
//存在更好的实现方法，比如直接把当前帧给清空用于执行
func tailCall(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1
	c := 0

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

//R(A+1) := R(B); R(A) := R(B)[RK(C)]
func self(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}
