package state

import "github.com/Youngkingman/GluaVirtual/luaState/luaApi"

//无法预知表的用法（数组or哈希）以及容量，使用NewTable创建
func (st *LuaState) NewTable() {
	st.CreateTable(0, 0)
}

func (st *LuaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	st.stack.push(t)
}

//根据栈顶值对索引的表进行查询
func (st *LuaState) GetTable(idx int) luaApi.LuaType {
	t := st.stack.get(idx)
	k := st.stack.pop()
	return st.getTable(t, k)
}

//根据给定的字符串对给定索引的表进行查询
func (st *LuaState) GetField(idx int, k string) luaApi.LuaType {
	t := st.stack.get(idx)
	return st.getTable(t, k)
}

//根据给定的i对给定索引的表（此时内部为数组）进行查询
func (st *LuaState) GetI(idx int, i int64) luaApi.LuaType {
	t := st.stack.get(idx)
	return st.getTable(t, i)
}

func (st *LuaState) getTable(t, k luaValue) luaApi.LuaType {
	if tab, ok := t.(*luaTable); ok {
		v := tab.get(k)
		st.stack.push(v)
		return typeOf(v)
	}
	panic("not a table")
}

func (st *LuaState) GetGlobal(name string) luaApi.LuaType {
	t := st.registry.get(luaApi.LUA_RIDX_GLOBALS) //获取注册表
	return st.getTable(t, name)                   //设置注册表
}

//将栈顶给出的值和键插入索引的表中
func (st *LuaState) SetTable(idx int) {
	t := st.stack.get(idx)
	v := st.stack.pop()
	k := st.stack.pop()
	st.setTable(t, k, v)
}

//对指定索引的表和字符串索引插入栈顶内容
func (st *LuaState) SetField(idx int, k string) {
	t := st.stack.get(idx)
	v := st.stack.pop()
	st.setTable(t, k, v)
}

//对指定索引的表（作为数组）插入栈顶内容
func (st *LuaState) SetI(idx int, n int64) {
	t := st.stack.get(idx)
	v := st.stack.pop()
	st.setTable(t, n, v)
}

func (st *LuaState) setTable(t, k, v luaValue) {
	if tab, ok := t.(*luaTable); ok {
		tab.put(k, v)
		return
	}
	panic("not a table")
}

func (st *LuaState) SetGlobal(name string) {
	t := st.registry.get(luaApi.LUA_RIDX_GLOBALS)
	v := st.stack.pop()
	st.setTable(t, name, v)
}

func (st *LuaState) Register(name string, f luaApi.GoFunction) {
	st.PushGoFunction(f)
	st.SetGlobal(name)
}
