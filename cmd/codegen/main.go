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
	funcs := make([]lgn.ExportFunc, 0)
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
		funcs = lgn.GenerateParams(funcs, parsed, fset, fileBytes)
	}
	fmt.Printf("Funcs size is %d\n", len(funcs))
	//fmt.Println(generateSource(pkgName, funcs))

	err = ioutil.WriteFile("statpkg_gen.go", lgn.GenerateSource(pkgName, funcs), 0644)
	//generateSource(pkgName, funcs)
	if err != nil {
		log.Fatal(err)
	}
}
