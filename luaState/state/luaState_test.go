package state

import (
	"fmt"
	"testing"

	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

func Test_ParseFunc(t *testing.T) {
	st := New()
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
