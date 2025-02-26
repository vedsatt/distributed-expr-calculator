package main

import (
	"fmt"

	"github.com/vedsatt/calc_prl/pkg/ast"
)

func main() {
	a := "1 + (2 - 3) * +4"
	_, err := ast.Build(a)
	if err != nil {
		fmt.Println(err)
	}
}
