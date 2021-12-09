package luaApi

//basic type from lua => go
type LuaType = int

const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

//stack arith operation type
type ArithOp = int

const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // -
	LUA_OPBNOT        // ~
)

//compare operation type
type CompareOp = int

const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)

type GoFunction func(LuaStateInterface) int

const LUA_MINSTACK = 20
const LUAI_MAXSTACK = 1000000                   //正负一百万为lua栈有效索引，正常够用
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000 //负的一百万是有效索引，减去1000为伪索引
const LUA_RIDX_GLOBALS int64 = 2

type LuaStateInterface interface {
	/* basic stack manipulation */
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(from, to int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)
	/* access functions (stack -> Go) */
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	IsTable(idx int) bool
	IsThread(idx int) bool
	IsFunction(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	/* push functions (Go -> stack) */
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	/*Arith Methods*/
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool
	Len(idx int)
	Concat(n int)
	/*lua table method*/
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, n int64)
	/*lua function call method*/
	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)
	/*go closure call method*/
	PushGoFunction(f GoFunction)
	IsGoFunction(idx int) bool
	ToGoFunction(idx int) GoFunction
	/*global environment support*/
	PushGlobalTable()
	GetGlobal(name string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)
	/*like C's pushCClosure*/
	PushGoClosure(f GoFunction, n int)
	/*manipulation of metatable and metamethod*/
	GetMetaTable(idx int) bool
	SetMetaTable(idx int)
	RawLen(idx int) uint
	RawEqual(idx1, idx2 int) bool
	RawGet(idx int) LuaType
	RawSet(idx int)
	RawGetI(idx int, i int64) LuaType
	RawSetI(idx int, i int64)
}

func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}
