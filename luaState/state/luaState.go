package state

import (
	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

var _ luaApi.LuaStateInterface = (*LuaState)(nil) //check implement of official luaApi
var _ luaApi.LuaVMInterface = (*LuaState)(nil)    //check extent luaApi for VM

type LuaState struct {
	stack *luaStack
	pc    int
	proto *binarychunk.Prototype
}

func New(stackSize int, proto *binarychunk.Prototype) *LuaState {
	return &LuaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
