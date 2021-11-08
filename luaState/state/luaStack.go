package state

type luaStack struct {
	slots []luaValue
	top   int
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
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
	if idx > 0 {
		return idx
	}
	return stk.top + idx + 1
}

func (stk *luaStack) isValid(idx int) bool {
	aIdx := stk.absIndex(idx)
	return aIdx > 0 && aIdx <= stk.top
}

func (stk *luaStack) get(idx int) luaValue {
	if stk.isValid(idx) {
		return stk.slots[stk.absIndex(idx)-1]
	}
	return nil
}

func (stk *luaStack) set(idx int, val luaValue) {
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
