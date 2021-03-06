package state

import (
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

var _ luaApi.LuaStateInterface = (*LuaState)(nil) //check implement of official luaApi
var _ luaApi.LuaVMInterface = (*LuaState)(nil)    //check extent luaApi for VM

type LuaState struct {
	registry *luaTable //状态机注册表
	stack    *luaStack
}

func New() *LuaState {
	registry := newLuaTable(0, 0)
	registry.put(luaApi.LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	st := &LuaState{registry: registry}
	st.pushLuaStack(newLuaStack(luaApi.LUA_MINSTACK, st))
	return st
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

func (st *LuaState) PushGoClosure(f luaApi.GoFunction, n int) {
	closure := newGoClosure(f, n)
	//将栈顶的n个数作为upvalue加入go的闭包中并传递
	for i := n; i > 0; i-- {
		val := st.stack.pop()
		closure.upvals[n-1] = &upvalue{&val}
	}
	st.stack.push(closure)
}
