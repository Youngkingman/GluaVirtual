package state

import (
	"math"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
	"github.com/Youngkingman/GluaVirtual/numTrans"
)

//table driven method
type operator struct {
	metamethod  string //对应操作的元方法
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
	{"__add", iadd, fadd},
	{"__sub", isub, fsub},
	{"__mul", imul, fmul},
	{"__mod", imod, fmod},
	{"__pow", nil, pow},
	{"__div", nil, div},
	{"__idiv", iidiv, fidiv},
	{"__band", band, nil},
	{"__bor", bor, nil},
	{"__bxor", bxor, nil},
	{"__shl", shl, nil},
	{"__shr", shr, nil},
	{"__unm", iunm, funm},
	{"__bnot", bnot, nil},
}

func (st *LuaState) Arith(op ArithOp) {
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
		return
	}

	metaMethod := operator.metamethod
	if result, ok := callMetamethod(a, b, metaMethod, st); ok {
		//操作数不可转换为数字时执行元方法的查找
		st.stack.push(result)
		return
	}
	panic("arithmetic error!")
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
