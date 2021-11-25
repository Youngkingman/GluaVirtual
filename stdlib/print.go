package stdlib

import (
	"fmt"

	"github.com/Youngkingman/GluaVirtual/luaState/luaApi"
)

func Print(st luaApi.LuaStateInterface) int {
	nArgs := st.GetTop()
	for i := 1; i <= nArgs; i++ {
		if st.IsBoolean(i) {
			fmt.Printf("%t", st.ToBoolean(i))
		} else {
			fmt.Print(st.TypeName(st.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}
