package main

import (
	"bytes"
	"fmt"
)

func generateSource(pkgName string, funcs []exportFunc) string {
	start := fmt.Sprintf("package main\nimport \"gonum.org/v1/gonum/%s\"\n\n", pkgName)
	buf := bytes.NewBufferString(start)
	for _, f := range funcs {
		buf.WriteString(fmt.Sprintf("%s%s(L *lua.LState) int {\n", pkgName, f.name))
		params := ""
		for i, p := range f.params {
			buf.WriteString(fmt.Sprintf("    x%d := %s(L, %d)\n", i+1, p.convertFunction(), i+1))
			if i > 0 {
				params = params + ", "
			}
			params = params + fmt.Sprintf("x%d", i)
		}
		rbit := ""
		for i := range f.returns {
			if i > 0 {
				rbit = rbit + ", "
			}
			rbit = rbit + fmt.Sprintf("r%d", i)
		}
		buf.WriteString(fmt.Sprintf("    %s := %s.%s(%s)\n", rbit, pkgName, f.name, params))
		for i, r := range f.returns {
			buf.WriteString(fmt.Sprintf("    %s(r%d)\n", convertReturnFunction(r), i))
		}
		buf.WriteString("    return 1\n}\n")
	}
	return buf.String()
}
