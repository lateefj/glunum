package main

import (
	"bytes"
	"fmt"

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
	lx := L.CheckTable(1)

	x := make([]float64, lx.Len())
	for i := 0; i < lx.Len(); i++ {
		if gv, ok := lx.RawGetInt(i).(lua.LNumber); ok {
			x[i] = float64(gv)
		}
	}
	nilWeights := L.Get(2)
	var w []float64
	if nilWeights != lua.LNil {
		lw := L.CheckTable(2)
		w = make([]float64, lw.Len())
		for i := 0; i < lw.Len(); i++ {
			if gv, ok := lw.RawGetInt(i).(lua.LNumber); ok {
				w[i] = float64(gv)
			}
		}
	}
	return x, w
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
}
