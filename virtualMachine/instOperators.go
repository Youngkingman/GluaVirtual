package vm

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

func _binaryArith(i Instruction, vm luaApi.LuaVMInterface, op luaApi.ArithOp) {
	a, b, c := i.ABC()
	a += 1
	//a寄存器用于存放结果，b/c可能是常量可能是寄存器里面的值，用GetRK堆到栈顶
	vm.GetRK(b)
	vm.GetRK(c)
	//进行二元计算
	vm.Arith(op)
	//存储结果
	vm.Replace(a)
}

func _unaryArith(i Instruction, vm luaApi.LuaVMInterface, op luaApi.ArithOp) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	//对b中的数进行指定位运算，结果存入a之中
	vm.GetRK(b)
	vm.Arith(op)
	vm.Replace(a)
}

//implement of binary or unary operations
func add(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPADD) }  // +
func sub(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPSUB) }  // -
func mul(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPMUL) }  // *
func mod(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPMOD) }  // %
func pow(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPPOW) }  // ^
func div(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPDIV) }  // /
func idiv(i Instruction, vm luaApi.LuaVMInterface) { _binaryArith(i, vm, luaApi.LUA_OPIDIV) } // //
func band(i Instruction, vm luaApi.LuaVMInterface) { _binaryArith(i, vm, luaApi.LUA_OPBAND) } // &
func bor(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPBOR) }  // |
func bxor(i Instruction, vm luaApi.LuaVMInterface) { _binaryArith(i, vm, luaApi.LUA_OPBXOR) } // ~
func shl(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPSHL) }  // <<
func shr(i Instruction, vm luaApi.LuaVMInterface)  { _binaryArith(i, vm, luaApi.LUA_OPSHR) }  // >>
func unm(i Instruction, vm luaApi.LuaVMInterface)  { _unaryArith(i, vm, luaApi.LUA_OPUNM) }   // -
func bnot(i Instruction, vm luaApi.LuaVMInterface) { _unaryArith(i, vm, luaApi.LUA_OPBNOT) }  // ~

func _len(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Len(b)
	vm.Replace(a)
}

func concat(i Instruction, vm luaApi.LuaVMInterface) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	c += 1
	//从b寄存器到c寄存器的值拼接，然后放入a中
	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	//栈顶n个拼接后值放入栈顶
	vm.Concat(n)
	vm.Replace(a)
}
