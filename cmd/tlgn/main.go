package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/lateefj/lgn"
)


func main() {

in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	l := lgn.NewLuaLoader(in, out, "./lua")
	statMod := l.State.SetFuncs(l.State.NewTable(), lgn.statFunctions)
	l.SetGlobalVar("stat", statMod)
	//statType := l.State.NewTypeMetatable(luaStatPackage)
	//l.SetGlobalVar("stat", statType)
	//l.SetField(statType, "__index", l.State.SetFuncs(l.State.NewTable(), statFunctions))
	err := l.File("idea.lua")
	if err != nil {
		fmt.Printf("Error %s", err)
	}

}
