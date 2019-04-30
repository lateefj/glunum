package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"gonum.org/v1/gonum/stat"
)

const (
	luaStatPackage = "stat"
)

var statFunctions = map[string]lua.LGFunction{
	"StdDev":   statStdDev,
	"Mean":     statMean,
	"Skew":     statSkew,
	"Variance": statVariance,
}

func twoTableExtract(L *lua.LState) ([]float64, []float64) {
	return paramFloatSlice(L, 1), paramFloatSlice(L, 2)
}
func statStdDev(L *lua.LState) int {
	x, w := twoTableExtract(L)
	resp := stat.StdDev(x, w)
	L.Push(lua.LNumber(resp))
	return 1
}

func statMean(L *lua.LState) int {
	x, w := twoTableExtract(L)
	resp := stat.Mean(x, w)
	L.Push(lua.LNumber(resp))
	return 1
}
func statVariance(L *lua.LState) int {
	x, w := twoTableExtract(L)
	resp := stat.Variance(x, w)
	L.Push(lua.LNumber(resp))
	return 1
}
func statSkew(L *lua.LState) int {
	x, w := twoTableExtract(L)
	resp := stat.Skew(x, w)
	L.Push(lua.LNumber(resp))
	return 1
}

// param ... Structure that keeps the name and type of parameter
type param struct {
	name  string
	ptype string
}

func (p *param) convertFunction() string {
	n, exists := paramConversionName[p.ptype]
	if n == "" || !exists {
		fmt.Println("************************************")
		fmt.Println("************Param Function Type*********************")
		fmt.Printf("%s\n", p.ptype)
		fmt.Println("************************************")
	}
	return n
}

func convertReturnFunction(rtype string) string {
	n, exists := returnConversionName[rtype]
	if n == "" || !exists {
		fmt.Println("************************************")
		fmt.Println("************Return Function Type*********************")
		fmt.Printf("%s\n", rtype)
		fmt.Println("************************************")
	}
	return n
}

// exportFunc ... Structure that stores the data needed to generate lua wrappers
type exportFunc struct {
	name    string
	params  []param
	returns []string
}

func generateParams(funcs []exportFunc, parsed *ast.File, fset *token.FileSet, fileBytes []byte) []exportFunc {
	ast.Inspect(parsed, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			if fn.Name.IsExported() {

				if fn.Type == nil || fn.Type.Params.List == nil || fn.Name.Obj == nil || fn.Name.Obj.Kind.String() != "func" {
					return true
				}
				params := make([]param, 0)
				for _, p := range fn.Type.Params.List {
					pos := fset.Position(p.Type.Pos())
					start := pos.Offset
					end := start + (int(p.Type.End()) - int(p.Type.Pos()))
					stype := string(fileBytes[start:end])
					params = append(params, param{name: p.Names[0].String(), ptype: stype})
				}
				results := make([]string, 0)
				if fn.Type.Results != nil {
					for _, r := range fn.Type.Results.List {
						pos := fset.Position(r.Type.Pos())
						start := pos.Offset
						end := start + (int(r.Type.End()) - int(r.Type.Pos()))
						stype := string(fileBytes[start:end])
						results = append(results, stype)
					}
				}
				funcs = append(funcs, exportFunc{name: fn.Name.Name, params: params, returns: results})
			}
			return true
		}
		return true
	})
	return funcs
}

func main() {

	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	l := NewLuaLoader(in, out, "./lua")
	//l.SetGlobalVar("std_dev", stat.StdDev)
	statMod := l.State.SetFuncs(l.State.NewTable(), statFunctions)
	l.SetGlobalVar("stat", statMod)
	/*statType := l.State.NewTypeMetatable(luaStatPackage)
	l.SetGlobalVar("stat", statType)
	l.SetField(statType, "__index", l.State.SetFuncs(l.State.NewTable(), statFunctions))*/
	err := l.File("idea.lua")
	if err != nil {
		fmt.Printf("Error %s", err)
	}

	pkgName := "stat"
	base := "./gonum/stat/"
	fset := token.NewFileSet()
	files, err := ioutil.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	funcs := make([]exportFunc, 0)
	for _, f := range files {
		if !strings.Contains(f.Name(), ".go") {
			continue
		}

		fmt.Println(f.Name())
		filePath := base + f.Name()
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			continue
		}
		parsed, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing: %s\n", err)
			continue
		}
		// Skip test files
		if strings.Contains(f.Name(), "_test.go") {
			continue
		}
		funcs = generateParams(funcs, parsed, fset, fileBytes)
	}
	fmt.Printf("Funcs size is %d\n", len(funcs))
	//fmt.Println(generateSource(pkgName, funcs))

	err = ioutil.WriteFile("statpkg_gen.go", generateSource(pkgName, funcs), 0644)
	//generateSource(pkgName, funcs)
	if err != nil {
		log.Fatal(err)
	}
}
