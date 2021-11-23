package state

import (
	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

var _ luaApi.LuaStateInterface = (*LuaState)(nil) //check implement of official luaApi
var _ luaApi.LuaVMInterface = (*LuaState)(nil)    //check extent luaApi for VM

type LuaState struct {
	stack *luaStack
}

func New(stackSize int, proto *binarychunk.Prototype) *LuaState {
	return &LuaState{
		stack: newLuaStack(stackSize),
	}
}

func (st *LuaState) pushLuaStack(stack *luaStack) {
	stack.prev = st.stack
	st.stack = stack
}

func (st *LuaState) popLuaStack() {
	stack := st.stack
	st.stack = stack.prev
	stack.prev = nil
}
