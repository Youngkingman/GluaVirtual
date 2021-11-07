package binarychunk

import "encoding/binary"

type reader struct {
	data []byte
}

//read one byte from byte stream
func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

//read bytes from byte stream
func (r *reader) readBytes(x uint) []byte {
	bytes := r.data[:x]
	r.data = r.data[x:]
	return bytes
}

//using little endian read uint32 from byte stream
func (r *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

//using little endian read uint64 from byte stream
func (r *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

//the luaInteger bytes map to int64 in go
func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

//the luaInteger bytes map to float64 in go
func (r *reader) readLuaNumber() float64 {
	return float64(r.readUint64())
}

//read string from byte stream
func (r *reader) readString() string {
	size := uint(r.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(r.readUint64())
	}
	bytes := r.readBytes(size)
	return string(bytes)
}

//see prototyStruct.PNG, its a recursion call, for function
//prototype is a recursion method
func (r *reader) readProto(source string) *Prototype {
	sourceName := r.readString()
	if sourceName == "" {
		sourceName = source
	}
	return &Prototype{
		Source:         sourceName,
		LineDefined:    r.readUint32(),
		LastLineDefine: r.readUint32(),
		NumParams:      r.readByte(),
		IsVararg:       r.readByte(),
		MaxStackSize:   r.readByte(),
		Code:           r.readCode(),
		Constants:      r.readSomeConstants(),
		Upvalues:       r.readUpvalues(),
		Protos:         r.readSomeProto(source),
		LineInfo:       r.readLineInfo(),
		LocVars:        r.readLocVars(),
		UpvalueNames:   r.readUpvalueNames(),
	}
}

//read code table from byte stream
func (r *reader) readCode() []uint32 {
	codes := make([]uint32, r.readUint32())
	for i := range codes {
		codes[i] = r.readUint32()
	}
	return codes
}

//read constant  from byte stream
func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() == 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STR:
		return r.readString()
	case TAG_LONG_STR:
		return r.readString()
	default:
		panic("read constant fail, corrupted!")
	}
}

//read constant table from byte streams
func (r *reader) readSomeConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

//read upvalues from byte stream
func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}

//read one proto from protos
func (r *reader) readSomeProto(source string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(source)
	}
	return protos
}

//see prototyStruct.PNG
func (r *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}
	return lineInfo
}

//see prototyStruct.PNG
func (r *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locVars
}

//see prototyStruct.PNG
func (r *reader) readUpvalueNames() []string {
	names := make([]string, r.readUint32())
	for i := range names {
		names[i] = r.readString()
	}
	return names
}

//check the chunk bytes header, following content in const&header
func (r *reader) checkHeader() {
	if string(r.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}
	if r.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if r.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(r.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if r.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}
	if r.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	}
	if r.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if r.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if r.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if r.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if r.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}
