package main

import (
	"github.com/vedsatt/calc_prl/internal/orchestrator"
)

func main() {
	app := orchestrator.New()

	app.Run()
}
