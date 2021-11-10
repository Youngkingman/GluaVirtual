package state

import (
	"fmt"
	"testing"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

func Test_StackBasicFunc(t *testing.T) {
	st := New(20, nil)
	st.PushBoolean(true)
	printStack(st)
	st.PushNumber(10)
	printStack(st)
	st.PushNil()
	printStack(st)
	st.PushString("fuck you")
	printStack(st)
	st.PushValue(-4)
	printStack(st)
	st.Replace(3)
	printStack(st)
	st.SetTop(6)
	printStack(st)
	st.Remove(-3)
	printStack(st)
	st.SetTop(-5)
	printStack(st)
}

func Test_StackArithMethod(t *testing.T) {
	st := New(20, nil)
	st.PushInteger(1)
	st.PushString("2.00")
	st.PushString("3.0")
	st.PushNumber(4.00)
	printStack(st)
	st.Arith(LUA_OPADD) //add 3.0 and 4.00
	printStack(st)
	st.Arith(LUA_OPBNOT) //reverse 7(0x0000000111)-> -8补码
	printStack(st)
	st.Len(2) //len of ["2.00"]
	printStack(st)
	st.Concat(3) //concat string from 1 to 3
	printStack(st)
	st.PushBoolean(st.Compare(1, 2, LUA_OPEQ))
	printStack(st)
}

func printStack(ls *LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}
