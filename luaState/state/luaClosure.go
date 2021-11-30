package state

import (
	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

//if proto is nil, then it's a goclosure
type luaClosure struct {
	proto  *binarychunk.Prototype
	goFunc luaApi.GoFunction
	upvals []*upvalue
}

type upvalue struct {
	val *luaValue
}

func newLuaClosure(proto *binarychunk.Prototype) *luaClosure {
	c := &luaClosure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(f luaApi.GoFunction, nUpvals int) *luaClosure {
	c := &luaClosure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
