package agent

import (
	"log"

	"github.com/vedsatt/calc_prl/internal/config"
)

type Agent struct {
	config config.Config
}

func New(cfg config.Config) *Agent {
	// передаем конфиг с переменными средами в агента
	return &Agent{config: cfg}
}

func (a *Agent) Run() {
	for i := 0; i < a.config.ComputingPower; i++ {
		log.Printf("worker %d starting", i+1)
		go worker(a.config)
	}

	select {} // бесконечное ожидание, решил без WaitGroup делать
}
