package main

import (
	lua "github.com/yuin/gopher-lua"
	//	"gonum.org/v1/gonum/mat"
)

var (
	paramConversionMap = map[string]interface{}{
		"[]float64": paramFloatArray,
		"float64":   paramFloat,
	}
	paramConversionName = map[string]string{
		"[]float64": "paramFloatArray",
		"float64":   "paramFloat",
		"int":       "paramInt",
		"[]bool":    "paramBoolArray",
		"bool":      "paramBool",
	}

	returnConversionMap = map[string]interface{}{
		"float64": returnFloat,
	}
	returnConversionName = map[string]string{
		"float64": "returnFloat",
	}
)

func returnFloat(L *lua.LState, v float64) {
	L.Push(lua.LNumber(v))
}

func paramBool(L *lua.LState, paramNumber int) bool {
	return bool(L.CheckBool(paramNumber))
}

func paramBoolArray(L *lua.LState, paramNumber int) []bool {
	nilCheck := L.Get(paramNumber)
	if nilCheck == lua.LNil {
		return nil
	}
	lx := L.CheckTable(paramNumber)

	x := make([]bool, lx.Len())
	for i := 0; i < lx.Len(); i++ {
		if gv, ok := lx.RawGetInt(i).(lua.LBool); ok {
			x[i] = bool(gv)
		}
	}
	return x
}
func paramInt(L *lua.LState, paramNumber int) int {
	return int(L.CheckInt(paramNumber))
}

func paramFloat(L *lua.LState, paramNumber int) float64 {
	return float64(L.CheckNumber(paramNumber))
}

func paramFloatArray(L *lua.LState, paramNumber int) []float64 {
	nilCheck := L.Get(paramNumber)
	if nilCheck == lua.LNil {
		return nil
	}
	lx := L.CheckTable(paramNumber)

	x := make([]float64, lx.Len())
	for i := 0; i < lx.Len(); i++ {
		if gv, ok := lx.RawGetInt(i).(lua.LNumber); ok {
			x[i] = float64(gv)
		}
	}
	return x
}

/*func returnDense(L *lua.LState, d *mat.Dense) {
}*/
