package state

import (
	"github.com/Youngkingman/GluaVirtual/binarychunk"
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

//if proto is nil, then it's a goclosure
type luaClosure struct {
	proto  *binarychunk.Prototype
	goFunc luaApi.GoFunction
}

func newLuaClosure(proto *binarychunk.Prototype) *luaClosure {
	return &luaClosure{proto: proto}
}

func newGoClosure(f luaApi.GoFunction) *luaClosure {
	return &luaClosure{goFunc: f}
}
