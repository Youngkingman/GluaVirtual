package state

func (st *LuaState) Len(idx int) {
	val := st.stack.get(idx)

	if s, ok := val.(string); ok {
		st.stack.push(int64(len(s)))
	} else if t, ok := val.(*luaTable); ok {
		st.stack.push(int64(t.tablen()))
	} else {
		panic("length error!")
	}
}

func (st *LuaState) Concat(n int) {
	if n == 0 {
		st.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if st.IsString(-1) && st.IsString(-2) {
				s2 := st.ToString(-1)
				s1 := st.ToString(-2)
				st.stack.pop()
				st.stack.pop()
				st.stack.push(s1 + s2)
				continue
			}
			panic("concatenation error!")
		}
	}
}
