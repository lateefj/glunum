package main

import (
	"bytes"
	"fmt"
)

func generateSource(pkgName string, funcs []exportFunc) string {
	start := fmt.Sprintf("package main\n import \"gonum.org/v1/gonum/%s\"\n", pkgName)
	buf := bytes.NewBufferString(start)
	for _, f := range funcs {
		buf.WriteString(fmt.Sprintf("%s%s(L *lua.LState) int {\n", pkgName, f.name))
		params := ""
		for i, p := range f.params {
			buf.WriteString(fmt.Sprintf("    x%d := %s(L, %d)\n", i+1, p.convertFunction(), i+1))
			params = params + fmt.Sprintf("")
		}
		buf.WriteString("}\n")
	}
	return buf.String()
}
