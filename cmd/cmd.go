package main

import "github.com/vedsatt/calc_prl/pkg/ast"

func main() {
	a := "2 + 3 * (4 - 1)"
	ast.Build(a)
}
