package main

import (
	lua "github.com/yuin/gopher-lua"
)

var (
	paramConversionMap = map[string]interface{}{
		"[]float64": paramFloatArray,
		"float64":   paramFloat,
	}

	returnConversionMap = map[string]interface{}{
		"float64": returnFloat,
	}
)

func returnFloat(L *lua.LState, v float64) {
	L.Push(lua.LNumber(v))
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
