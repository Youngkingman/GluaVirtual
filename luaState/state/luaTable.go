package state

import "github.com/Youngkingman/GluaVirtual/numTrans"

type luaTable struct {
	luaArr []luaValue
	luaMap map[luaValue]luaValue
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.luaArr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t.luaMap = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (tab *luaTable) get(key luaValue) luaValue {
	key = _float2Integer(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(tab.luaArr)) {
			return tab.luaArr[idx-1]
		}
	}
	return tab.luaMap[key]
}

func _float2Integer(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := numTrans.Float2Integer(f); ok {
			return i
		}
	}
	return key
}
