package luaApi

type LuaVMInterface interface {
	LuaStateInterface
	PC() int          //get pc information, only for test
	OffsetPC(n int)   //used in jump istr to modify pc
	Fetch() uint32    //fetch current inctruction, move to next istr
	GetConst(idx int) //get const from constants table of proto, and push it into luaState
	GetRK(rk int)
}
