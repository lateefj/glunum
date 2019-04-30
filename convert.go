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
		"mat.Vector":    "paramVector",
		"*mat.Cholesky": "paramCholeskyPointer",
	}

	returnConversionMap = map[string]interface{}{
		"float64": returnFloat,
	}
	returnConversionName = map[string]string{
		"[]float64":     "returnFloatSlice",
		"float64":       "returnFloat",
		"mat.Matrix":    "returnMatrix",
		"*mat.SymDense": "returnSymDensePointer",
		"mat.Vector":    "returnVector",
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

func vectorFuncs(tbl *lua.LTable, v mat.Vector) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"Dims": func(L *lua.LState) int {
			r, c := v.Dims()
			L.Push(lua.LNumber(r))
			L.Push(lua.LNumber(c))
			return 1
		},
		"At": func(L *lua.LState) int {
			i := paramInt(L, 1)
			j := paramInt(L, 2)
			L.Push(lua.LNumber(v.At(i, j)))
			return 1
		},
		"AtVec": func(L *lua.LState) int {
			i := paramInt(L, 1)
			L.Push(lua.LNumber(v.AtVec(i)))
			return 1
		},
		"Len": func(L *lua.LState) int {
			L.Push(lua.LNumber(v.Len()))
			return 1
		},
	}
}
func returnVector(L *lua.LState, v mat.Vector) *lua.LTable {
	tbl := L.NewTable()
	L.SetFuncs(tbl, vectorFuncs(tbl, v))
	L.Push(tbl)
	return tbl
}

func symDensePointerFuncs(tbl *lua.LTable, sd *mat.SymDense) map[string]lua.LGFunction {
	funcs := matrixFuncs(tbl, sd)
	funcs["Symmetric"] = func(L *lua.LState) int {
		L.Push(lua.LNumber(sd.Symmetric()))
		return 1
	}
	return funcs
}

func returnSynDensePointer(L *lua.LState, sd *mat.SymDense) *lua.LTable {
	tbl := L.NewTable()
	L.SetFuncs(tbl, symDensePointerFuncs(tbl, sd))
	L.Push(tbl)
	return tbl
}

func returnCholeskyPointer(L *lua.LState, c *mat.Cholesky) *lua.LTable {
	tbl := L.NewTable()
	L.SetFuncs(tbl, choleskyPointerFuncs(tbl, c))
	L.Push(tbl)
	return tbl
}

type wrapCholeskyPointer struct {
	L     *lua.LState
	table *lua.LTable
}

func (wcp *wrapCholeskyPointer) Cond() float64 {
	wcp.L.CallByParam(lua.P{
		Fn:      wcp.L.GetGlobal("Cond"),
		NRet:    1,
		Protect: true,
	},
	)
	x := float64(wcp.L.ToNumber(-1))
	wcp.L.Pop(1)
	return x
}
func (wcp *wrapCholeskyPointer) Det() float64 {
	wcp.L.CallByParam(lua.P{
		Fn:      wcp.L.GetGlobal("Det"),
		NRet:    1,
		Protect: true,
	},
	)
	x := float64(wcp.L.ToNumber(-1))
	wcp.L.Pop(1)
	return x
}
func (wcp *wrapCholeskyPointer) ExtendVecSym(a *mat.Cholesky, v mat.Vector) bool {
	wcp.L.CallByParam(lua.P{
		Fn:      wcp.L.GetGlobal("ExtendVecSym"),
		NRet:    1,
		Protect: true,
	}, returnCholeskyPointer(wcp.L, a), returnVector(wcp.L, v),
	)
	x := wcp.L.ToBool(-1)
	wcp.L.Pop(1)
	return x
}

func (wcp *wrapCholeskyPointer) Clone(a *mat.Cholesky) {
	//return &wrapCholeskyPointer{L: wcp.L, table: wcp.table}
}
func paramCholeskyPointer(L *lua.LState, paramNumber int) *mat.Cholesky {
	return &mat.Cholesky{}
}

func choleskyPointerFuncs(tbl *lua.LTable, c *mat.Cholesky) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"Clone": func(L *lua.LState) int {
			a := paramCholeskyPointer(L, 1)
			c.Clone(a)
			return 1
		},
		"Cond": func(L *lua.LState) int {
			L.Push(lua.LNumber(c.Cond()))
			return 1
		},
		"Det": func(L *lua.LState) int {
			L.Push(lua.LNumber(c.Det()))
			return 1
		},
		"ExtendVecSym": func(L *lua.LState) int {
			a := paramCholeskyPointer(L, 1)
			v := paramVector(L, 2)
			L.Push(lua.LBool(c.ExtendVecSym(a, v)))
			return 1
		}}
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
	x := wm.L.ToInt(-1)
	wm.L.Pop(1)
	y := wm.L.ToInt(-1)
	wm.L.Pop(1)
	return x, y
}

func (wm *wrapMatrix) At(i, j int) float64 {
	wm.L.CallByParam(lua.P{
		Fn:      wm.L.GetGlobal("At"),
		NRet:    1,
		Protect: true,
	}, lua.LNumber(i), lua.LNumber(j),
	)
	x := float64(wm.L.ToNumber(-1))
	wm.L.Pop(1)
	return x
}
func (wm *wrapMatrix) T() mat.Matrix {
	return &wrapMatrix{L: wm.L, table: wm.table}
}
func paramMatrix(L *lua.LState, paramNumber int) mat.Matrix {
	return &wrapMatrix{L: L, table: L.CheckTable(paramNumber)}
}

type wrapVector struct {
	wrapMatrix
}

func (wv wrapVector) AtVec(i int) float64 {
	wv.L.CallByParam(lua.P{
		Fn:      wv.L.GetGlobal("AtVec"),
		NRet:    1,
		Protect: true,
	}, lua.LNumber(i),
	)
	x := float64(wv.L.ToNumber(-1))
	wv.L.Pop(1)
	return x
}

func (wv *wrapVector) Len() int {
	wv.L.CallByParam(lua.P{
		Fn:      wv.L.GetGlobal("Len"),
		NRet:    1,
		Protect: true,
	},
	)
	x := wv.L.ToInt(-1)
	wv.L.Pop(1)
	return x
}

func paramVector(L *lua.LState, paramNumber int) mat.Vector {
	return &wrapVector{wrapMatrix{L: L, table: L.CheckTable(paramNumber)}}
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
