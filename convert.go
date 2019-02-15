package main

import (
	lua "github.com/yuin/gopher-lua"
	//	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

var (
	paramConversionMap = map[string]interface{}{
		"[]float64": paramFloatSlice,
		"float64":   paramFloat,
	}
	paramConversionName = map[string]string{
		"[]float64":     "paramFloatSlice",
		"float64":       "paramFloat",
		"int":           "paramInt",
		"[]bool":        "paramBoolSlice",
		"bool":          "paramBool",
		"CumulantKind":  "paramCumulantKind",
		"mat.Matrix":    "paramMatrix",
		"*mat.SymDense": "paramSymDensePointer",
	}

	returnConversionMap = map[string]interface{}{
		"float64": returnFloat,
	}
	returnConversionName = map[string]string{
		"[]float64":     "returnFloatSlice",
		"float64":       "returnFloat",
		"mat.Matrix":    "returnMatrix",
		"*mat.SymDense": "returnSymDensePointer",
	}
)

func returnFloat(L *lua.LState, v float64) {
	L.Push(lua.LNumber(v))
}

func matrixFuncs(tbl *lua.LTable, m mat.Matrix) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"Dims": func(L *lua.LState) int {
			r, c := m.Dims()
			L.Push(lua.LNumber(r))
			L.Push(lua.LNumber(c))
			return 1
		},
		"At": func(L *lua.LState) int {
			i := paramInt(L, 1)
			j := paramInt(L, 2)
			L.Push(lua.LNumber(m.At(i, j)))
			return 1
		},
	}
}

func returnMatrix(L *lua.LState, m mat.Matrix) {
	tbl := L.NewTable()
	L.SetFuncs(tbl, matrixFuncs(tbl, m))
	L.Push(tbl)
}

func symDensePointerFuncs(tbl *lua.LTable, sd *mat.SymDense) map[string]lua.LGFunction {
	funcs := matrixFuncs(tbl, sd)
	funcs["Symmetric"] = func(L *lua.LState) int {
		L.Push(lua.LNumber(sd.Symmetric()))
		return 1
	}
	return funcs
}

func returnSynDensePointer(L *lua.LState, sd *mat.SymDense) {
	tbl := L.NewTable()
	L.SetFuncs(tbl, symDensePointerFuncs(tbl, sd))
	L.Push(tbl)
}

type wrapMatrix struct {
	L     *lua.LState
	table *lua.LTable
}

func (wm *wrapMatrix) Dims() (int, int) {
	wm.L.CallByParam(lua.P{
		Fn:      wm.L.GetGlobal("Dims"),
		NRet:    2,
		Protect: true,
	},
	)
	x := wm.L.Get(-1)
	wm.L.Pop(1)
	y := wm.L.Get(-1)
	wm.L.Pop(1)
	return wm.L.CheckIn
}

func paramMatrix(L *lua.LState, paramNumber int) mat.Matrix {
	wm := wrapMatrix{L.CheckTable(paramNumber)}
	return wm
}

/*func paramSymDensePointer(L *lua.LState, paramNumber int) *mat.SymDense {
	tbl := L.CheckTable(paramNumber)
	tbl.

	L.SetFunc(tbl, symDensePointerFuncs)
	L.Push(tbl)
}*/

func paramCumulantKind(L *lua.LState, paramNumber int) stat.CumulantKind {
	return stat.CumulantKind(L.CheckInt(paramNumber))
}
func paramBool(L *lua.LState, paramNumber int) bool {
	return bool(L.CheckBool(paramNumber))
}

func paramBoolSlice(L *lua.LState, paramNumber int) []bool {
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

func paramFloatSlice(L *lua.LState, paramNumber int) []float64 {
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

func returnFloatSlice(L *lua.LState, s []float64) {
	tbl := L.NewTable()
	for i, v := range s {
		tbl.Insert(i, lua.LNumber(v))
	}
	L.Push(tbl)
}
