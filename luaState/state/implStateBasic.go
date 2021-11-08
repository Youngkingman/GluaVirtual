package state

import (
	. "github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

var _ LuaStateInterface = (*luaState)(nil)

func (st *luaState) GetTop() int {
	return st.stack.top
}

func (st *luaState) AbsIndex(idx int) int {
	return st.stack.absIndex(idx)
}

//didn't consider the situation of failure in expansion of memory
func (st *luaState) CheckStack(n int) bool {
	st.stack.check(n)
	return true
}

func (st *luaState) Pop(n int) {
	for n > 0 {
		st.stack.pop()
		n--
	}
}

//copy a value from certain index to certain index
func (st *luaState) Copy(from, to int) {
	val := st.stack.get(from)
	st.stack.set(to, val)
}

//push value in certain index into the stack
func (st *luaState) PushValue(idx int) {
	val := st.stack.get(idx)
	st.stack.push(val)
}

//inverse operation of PushValue
func (st *luaState) Replace(idx int) {
	val := st.stack.pop()
	st.stack.set(idx, val)
}

//typical leetcode problem,rotate n the array [absidx,top]
func (st *luaState) Rotate(idx, n int) {
	t := st.stack.top - 1
	p := st.stack.absIndex(idx) - 1
	m := -1
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	st.stack.inverse(p, m)
	st.stack.inverse(m+1, t)
	st.stack.inverse(p, t)
}

//pop top value and insert into certain index
func (st *luaState) Insert(idx int) {
	st.Rotate(idx, 1)
}

//remove value in certain index
func (st *luaState) Remove(idx int) {
	st.Rotate(idx, -1)
	st.stack.pop()
}

//let element in certain index(valid or overflow invalid) become top element
func (st *luaState) SetTop(idx int) {
	newTop := st.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow invalid index")
	}
	n := st.stack.top - newTop
	if n > 0 {
		st.Pop(n)
	} else {
		for i := 0; i > n; i-- {
			st.stack.push(nil)
		}
	}
}
