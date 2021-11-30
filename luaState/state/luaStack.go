package state

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

type luaStack struct {
	slots []luaValue
	top   int

	//function call stack segment
	prev    *luaStack
	closure *luaClosure
	varargs []luaValue
	pc      int
	//used to access registry
	state *LuaState
	//map of upvalues, open upvalue or closed upvalue
	openuvs map[int]*upvalue
}

func newLuaStack(size int, st *LuaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: st,
	}
}

//check if there is enough space, if not, expand the space
func (stk *luaStack) check(n int) {
	free := len(stk.slots) - stk.top
	for i := free; i < n; i++ {
		stk.slots = append(stk.slots, nil)
	}
}

func (stk *luaStack) push(val luaValue) {
	if stk.top == len(stk.slots) {
		panic("stack overflow")
	}
	stk.slots[stk.top] = val
	stk.top++
}

func (stk *luaStack) pop() (ret luaValue) {
	if stk.top < 1 {
		panic("stack is already empty")
	}
	stk.top--
	ret = stk.slots[stk.top]
	stk.slots[stk.top] = nil
	return
}

//tranfer the relavtive index into absolute index in the stack
func (stk *luaStack) absIndex(idx int) int {
	if idx <= luaApi.LUA_REGISTRYINDEX {
		return idx //伪索引，不用转换
	}
	if idx > 0 {
		return idx
	}
	return stk.top + idx + 1
}

func (stk *luaStack) isValid(idx int) bool {
	if idx < luaApi.LUA_REGISTRYINDEX {
		//小于注册表索引说明是upvalue的伪索引
		uvIdx := luaApi.LUA_REGISTRYINDEX - idx - 1
		c := stk.closure
		return c != nil && uvIdx < len(c.upvals)
	}
	if idx == luaApi.LUA_REGISTRYINDEX {
		//注册表有效
		return true
	}
	aIdx := stk.absIndex(idx)
	return aIdx > 0 && aIdx <= stk.top
}

func (stk *luaStack) get(idx int) luaValue {
	if idx < luaApi.LUA_REGISTRYINDEX {
		//小于注册表索引说明是upvalue的伪索引
		uvidx := luaApi.LUA_REGISTRYINDEX - idx - 1
		c := stk.closure
		if c == nil || uvidx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvidx].val)
	}
	if idx == luaApi.LUA_REGISTRYINDEX {
		//直接返回注册表
		return stk.state.registry
	}
	if stk.isValid(idx) {
		return stk.slots[stk.absIndex(idx)-1]
	}
	return nil
}

func (stk *luaStack) set(idx int, val luaValue) {
	if idx < luaApi.LUA_REGISTRYINDEX {
		//小于注册表索引说明是upvalue的伪索引
		uvidx := luaApi.LUA_REGISTRYINDEX - idx - 1
		c := stk.closure
		if c != nil && uvidx < len(c.upvals) {
			*(c.upvals[uvidx].val) = val
		}
		return
	}
	if idx == luaApi.LUA_REGISTRYINDEX {
		//注册表设置不用转换,如果传入val是nil可能会清空注册表
		stk.state.registry = val.(*luaTable)
		return
	}
	if stk.isValid(idx) {
		stk.slots[stk.absIndex(idx)-1] = val
		return
	}
	panic("invalid index")
}

func (stk *luaStack) inverse(from, to int) {
	for from < to {
		stk.slots[from], stk.slots[to] = stk.slots[to], stk.slots[from]
		from++
		to--
	}
}

func (stk *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = stk.pop()
	}
	return vals
}

func (stk *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			stk.push(vals[i])
		} else {
			stk.push(nil)
		}
	}
}
