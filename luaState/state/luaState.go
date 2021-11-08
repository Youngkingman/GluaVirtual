package state

type LuaState struct {
	stack *luaStack
}

func New() *LuaState {
	return &LuaState{
		stack: newLuaStack(20),
	}
}
