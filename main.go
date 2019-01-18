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
	return paramFloatArray(L, 1), paramFloatArray(L, 2)
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

	base := "./gonum/stat/"
	fset := token.NewFileSet()
	files, err := ioutil.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
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
		ast.Inspect(parsed, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if ok {
				if fn.Name.IsExported() {
					fmt.Printf("Function: %s\n", fn.Name.Name)

					params := []string{}
					if fn.Type == nil || fn.Type.Params.List == nil {
						return true
					}
					for _, p := range fn.Type.Params.List {
						fmt.Printf("%T\n", p.Type)
						//stype := fmt.Sprintf("%T", p.Type)
						//stype := fset.Position(p.Type.Pos()).String()
						/*if p.Tag == nil {
							continue
						}*/
						//stype := fmt.Sprintf("%T", p.Type)
						//stype := p.Tag.Kind.String()
						pos := fset.Position(p.Type.Pos())
						fmt.Printf("Pos : %d\n", pos.Offset)
						//stype := p.Names[0].Obj.Kind.String()
						start := pos.Offset
						end := start + (int(p.Type.End()) - int(p.Type.Pos()))
						fmt.Printf("start %d and end %d total bytes %d\n", start, end, len(fileBytes))
						stype := string(fileBytes[start:end])
						params = append(params, fmt.Sprintf("\ttype name: %s type: %T type: %+v src: %s\n", p.Names[0], p.Type, p.Type, stype))
					}
					fmt.Printf("\n\tParams: %s\n", strings.Join(params, ","))
				}
				return true
			}
			return true
		})
	}
}
