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
	return st.getTable(t, k, false) //false触发元方法，下同
}

//根据给定的字符串对给定索引的表进行查询
func (st *LuaState) GetField(idx int, k string) luaApi.LuaType {
	t := st.stack.get(idx)
	return st.getTable(t, k, false)
}

//根据给定的i对给定索引的表（此时内部为数组）进行查询
func (st *LuaState) GetI(idx int, i int64) luaApi.LuaType {
	t := st.stack.get(idx)
	return st.getTable(t, i, false)
}

func (st *LuaState) getTable(t, k luaValue, raw bool) luaApi.LuaType {
	if tab, ok := t.(*luaTable); ok {
		v := tab.get(k)
		//如果raw == true/v不存在/没有元索引字段 则忽略元方法
		if raw || v != nil || !tab.hasMetafield("__index") {
			st.stack.push(v)
			return typeOf(v)
		}
	}

	if !raw {
		if mf := getMetafield(t, "__index", st); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return st.getTable(x, k, false) //递归查找方法字段，直到找到元方法并执行
			case *luaClosure:
				//执行元方法
				st.stack.push(mf)
				st.stack.push(t)
				st.stack.push(k)
				st.Call(2, 1)
				v := st.stack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("index error")
}

func (st *LuaState) GetGlobal(name string) luaApi.LuaType {
	t := st.registry.get(luaApi.LUA_RIDX_GLOBALS) //获取注册表
	return st.getTable(t, name, false)            //设置注册表
}

//将栈顶给出的值和键插入索引的表中
func (st *LuaState) SetTable(idx int) {
	t := st.stack.get(idx)
	v := st.stack.pop()
	k := st.stack.pop()
	st.setTable(t, k, v, false)
}

//对指定索引的表和字符串索引插入栈顶内容
func (st *LuaState) SetField(idx int, k string) {
	t := st.stack.get(idx)
	v := st.stack.pop()
	st.setTable(t, k, v, false)
}

//对指定索引的表（作为数组）插入栈顶内容
func (st *LuaState) SetI(idx int, n int64) {
	t := st.stack.get(idx)
	v := st.stack.pop()
	st.setTable(t, n, v, false)
}

func (st *LuaState) setTable(t, k, v luaValue, raw bool) {
	if tab, ok := t.(*luaTable); ok {
		if raw || tab.get(k) != nil || !tab.hasMetafield("__newindex") {
			tab.put(k, v)
		}
		return
	}

	if !raw {
		if mf := getMetafield(t, "__newindex", st); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				st.setTable(x, k, v, false)
				return
			case *luaClosure:
				st.stack.push(mf)
				st.stack.push(t)
				st.stack.push(k)
				st.stack.push(v)
				st.Call(3, 0)
				return
			}
		}
	}
	panic("index error")
}

func (st *LuaState) SetGlobal(name string) {
	t := st.registry.get(luaApi.LUA_RIDX_GLOBALS)
	v := st.stack.pop()
	st.setTable(t, name, v, false)
}

func (st *LuaState) Register(name string, f luaApi.GoFunction) {
	st.PushGoFunction(f)
	st.SetGlobal(name)
}
