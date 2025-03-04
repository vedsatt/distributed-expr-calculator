package main

import (
	"github.com/vedsatt/calc_prl/internal/agent"
	"github.com/vedsatt/calc_prl/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	agent := agent.New(cfg)
	agent.Run()
}
