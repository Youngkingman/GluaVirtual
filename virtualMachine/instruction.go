package vm

type Instruction uint32 // Instruction

// BX  |------------------|------------------|
//     0                131071              262143
// SBX |------------------|------------------|
//  -131071               0                 131072

const MAXARG_Bx = 1<<18 - 1
const MAXARG_sBx = MAXARG_Bx >> 1

//GET operation code from single instruction
func (istr Instruction) Opcode() int {
	return int(istr & 0x3F)
}

//parse operation number from instruction in ABC mode
func (istr Instruction) ABC() (a, b, c int) {
	a = int(istr >> 6 & 0xFF)
	c = int(istr >> 14 & 0x1FF)
	b = int(istr >> 23 & 0x1FF)
	return
}

//parse operation number from instruction in ABx mode
func (istr Instruction) ABx() (a, bx int) {
	a = int(istr >> 6 & 0xFF)
	bx = int(istr >> 14 & 0x1FF)
	return
}

//parse operation number from instruction in AsBx mode
//sbx is a signed number encoded in offset binary
func (istr Instruction) AsBx() (a, sbx int) {
	a, sbx = istr.ABx()
	sbx -= MAXARG_sBx
	return
}

//parse operation number from instruction in Ax mode
func (istr Instruction) Ax() int {
	return int(istr >> 6)
}

//get name of the instruction
func (istr Instruction) OpName() string {
	return opcodes[istr.Opcode()].name
}

//get mode of the instruction
func (istr Instruction) OpMode() byte {
	return opcodes[istr.Opcode()].opMode
}

//get ArgBMode from instruction
func (istr Instruction) ArgBMode() byte {
	return opcodes[istr.Opcode()].argBMode
}

//get argCMode from instruction
func (istr Instruction) ArgCMode() byte {
	return opcodes[istr.Opcode()].argCMode
}
