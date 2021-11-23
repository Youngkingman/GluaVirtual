package state

import "github.com/Youngkingman/GluaVirtual/binarychunk"

type luaClosure struct {
	proto *binarychunk.Prototype
}

func newLuaClosure(proto *binarychunk.Prototype) *luaClosure {
	return &luaClosure{proto: proto}
}
