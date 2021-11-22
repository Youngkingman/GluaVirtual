package vm

/*
	将一个整数转换为浮点字节，(eeeeexxx)当eeeee==0时字节表示的整数是xxx，
	否则这个字节表示的是(1xxx) * 2^(eeeee-1)
	主要用于初始化表的长度,正常操作数(B和C只有9个比特，最大也就索引512个，采用浮点数
	可以有最大15^(32-1)的范围表示(当然没法对应每一个数），用于JSON等对象初始化时不需要每次都进行扩容

*/
func Int2fb(x int) int {
	e := 0 /* exponent */
	if x < 8 {
		return x
	}
	for x >= (8 << 4) { /* coarse steps */
		x = (x + 0xf) >> 4 /* x = ceil(x / 16) */
		e += 4
	}
	for x >= (8 << 1) { /* fine steps */
		x = (x + 1) >> 1 /* x = ceil(x / 2) */
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

/* converts back */
func Fb2int(x int) int {
	if x < 8 {
		return x
	} else {
		return ((x & 7) + 8) << uint((x>>3)-1)
	}
}
