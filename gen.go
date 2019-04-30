package main

import (
	"bytes"
	"fmt"
)

func generateSource(pkgName string, funcs []exportFunc) []byte {
	start := fmt.Sprintf("package main\nimport (\"gonum.org/v1/gonum/%s\"\n\tlua \"github.com/yuin/gopher-lua\"\n)\n", pkgName)
	buf := bytes.NewBufferString(start)
	//  Add all function names to export func
	buf.WriteString(fmt.Sprintf("var %sFunctions = map[string]lua.LGFunction {\n", pkgName))
	for _, f := range funcs {
		buf.WriteString(fmt.Sprintf("\t\"%s\": %s%s,\n", f.name, pkgName, f.name))
	}
	buf.WriteString("}\n")
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
	return buf.Bytes()
}
