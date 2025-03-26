package agent

import (
	"log"
	"time"

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
	for i := range a.config.ComputingPower {
		time.Sleep(1 * time.Second)
		log.Printf("worker %d starting...", i+1)
		go worker(a.config)
	}

	select {} // бесконечное ожидание
}
