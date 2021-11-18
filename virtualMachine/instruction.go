package vm

import (
	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

type Instruction uint32 // Instruction

// BX  |------------------|------------------|
//     0                131071              262143
// SBX |------------------|------------------|
//  -131071               0                 131072

const MAXARG_Bx = 1<<18 - 1
const MAXARG_sBx = MAXARG_Bx >> 1

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/

//GET operation code from single instruction
func (inst Instruction) Opcode() int {
	return int(inst & 0x3F)
}

//parse operation number from instruction in ABC mode
func (inst Instruction) ABC() (a, b, c int) {
	a = int(inst >> 6 & 0xFF)
	c = int(inst >> 14 & 0x1FF)
	b = int(inst >> 23 & 0x1FF)
	return
}

//parse operation number from instruction in ABx mode
func (inst Instruction) ABx() (a, bx int) {
	a = int(inst >> 6 & 0xFF)
	bx = int(inst >> 14)
	return
}

//parse operation number from instruction in AsBx mode
//sbx is a signed number encoded in offset binary
func (inst Instruction) AsBx() (a, sbx int) {
	a, bx := inst.ABx()
	return a, bx - MAXARG_sBx
}

//parse operation number from instruction in Ax mode
func (inst Instruction) Ax() int {
	return int(inst >> 6)
}

//get name of the instruction
func (inst Instruction) OpName() string {
	return opcodes[inst.Opcode()].name
}

//get mode of the instruction
func (inst Instruction) OpMode() byte {
	return opcodes[inst.Opcode()].opMode
}

//get ArgBMode from instruction
func (inst Instruction) ArgBMode() byte {
	return opcodes[inst.Opcode()].argBMode
}

//get argCMode from instruction
func (inst Instruction) ArgCMode() byte {
	return opcodes[inst.Opcode()].argCMode
}

func (inst Instruction) Execute(vm luaApi.LuaVMInterface) {
	action := opcodes[inst.Opcode()].action
	if action != nil {
		action(inst, vm)
	} else {
		panic(inst.OpName() + " currently dosen't have action")
	}
}
