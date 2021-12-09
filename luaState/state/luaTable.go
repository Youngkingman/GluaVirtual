package state

import (
	"math"

	"github.com/Youngkingman/GluaVirtual/numTrans"
)

type luaTable struct {
	luaArr    []luaValue
	luaMap    map[luaValue]luaValue
	metatable *luaTable
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

func (tab *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("nil key is unacceptable")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("NaN key is unacceptable")
	}
	key = _float2Integer(key)

	//整数索引使用数组
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(tab.luaArr))
		if idx <= arrLen {
			tab.luaArr[idx-1] = val
			if idx == arrLen && val == nil {
				tab.shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(tab.luaMap, key)
			if val != nil {
				tab.luaArr = append(tab.luaArr, val)
				tab.expandArray()
			}
			return
		}
	}

	if val != nil {
		if tab.luaMap == nil {
			tab.luaMap = make(map[luaValue]luaValue, 8)
		}
		tab.luaMap[key] = val
	} else {
		delete(tab.luaMap, key)
	}
}

func (tab *luaTable) shrinkArray() {
	for i := len(tab.luaArr) - 1; i >= 0; i-- {
		if tab.luaArr[i] == nil {
			tab.luaArr = tab.luaArr[0:i]
		}
	}
}

func (tab *luaTable) expandArray() {
	for idx := int64(len(tab.luaArr)) + 1; true; idx++ {
		if val, has := tab.luaMap[idx]; has {
			delete(tab.luaMap, idx)
			tab.luaArr = append(tab.luaArr, val)
		} else {
			break
		}
	}
}

func (tab *luaTable) tablen() int {
	return len(tab.luaArr)
}

func (tab *luaTable) hasMetafield(fieldname string) bool {
	return tab.metatable != nil && tab.metatable.get(fieldname) != nil
}
