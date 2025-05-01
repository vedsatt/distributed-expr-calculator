package agent

import (
	"log"

	"github.com/vedsatt/calc_prl/internal/config"
	"github.com/vedsatt/calc_prl/internal/models"
)

type Agent struct {
	config config.Config
}

type Task struct {
	ID   int
	Arg1 string
	Arg2 string
	Type string
}

var (
	resultsCh = make(chan *models.Result)
	tasksCh   = make(chan *Task)
)

func New(cfg config.Config) *Agent {
	// передаем конфиг с переменными средами в агента
	return &Agent{config: cfg}
}

func (a *Agent) Run() {
	go Connect()

	for i := range a.config.ComputingPower {
		log.Printf("worker %d starting...", i+1)
		go worker(a.config)
	}

	select {} // бесконечное ожидание
}
