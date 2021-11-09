package state

import (
	"math"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
	"github.com/Youngkingman/GluaVirtual/numTrans"
)

//table driven method
type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = numTrans.IntegerMod
	fmod  = numTrans.FloatMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = numTrans.IntegerFloorDiv
	fidiv = numTrans.FloatFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = numTrans.LeftShift
	shr   = numTrans.RightShift
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	{iadd, fadd},
	{isub, fsub},
	{imul, fmul},
	{imod, fmod},
	{nil, pow},
	{nil, div},
	{iidiv, fidiv},
	{band, nil},
	{bor, nil},
	{bxor, nil},
	{shl, nil},
	{shr, nil},
	{iunm, funm},
	{bnot, nil},
}

func (st *LuaState) ArithOp(op ArithOp) {
	var a, b luaValue //operands
	b = st.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = st.stack.pop()
	} else {
		a = b
	}
	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		st.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil { //bit operation like and or xor shl shr not
		x, oka := convertToInteger(a)
		y, okb := convertToInteger(b)
		if oka && okb {
			return op.integerFunc(x, y)
		}
	} else {
		if op.integerFunc != nil { // add sub mul mod idiv unm
			x, oka := a.(int64)
			y, okb := b.(int64)
			if oka && okb {
				return op.integerFunc(x, y)
			}
		}
		x, oka := converToFloat(a)
		y, okb := converToFloat(b)
		if oka && okb {
			return op.floatFunc(x, y)
		}
	}
	return nil
}
