package luaApi

type LuaVMInterface interface {
	LuaStateInterface
	PC() int          //get pc information, only for test
	OffsetPC(n int)   //used in jump istr to modify pc
	Fetch() uint32    //fetch current inctruction, move to next istr
	GetConst(idx int) //get const from constants table of proto, and push it into luaState
	GetRK(rk int)     //RK may be a register or const, this method will return the value due to the input

	RegisterCount() int //return current function frame's register count
	LoadVararg(n int)   //load varargs to current function frame
	LoadProto(idx int)  //transfer a closure into function proto and push it into the stack
}
