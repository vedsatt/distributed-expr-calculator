package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vedsatt/calc_prl/internal/config"
	"github.com/vedsatt/calc_prl/pkg/ast"
)

type Result struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
	Error  string  `json:"error"`
}

func getTask() (*ast.AstNode, int) {
	client := &http.Client{
		Timeout: 5 * time.Second, // Таймаут 5 секунд
	}

	resp, err := client.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, http.StatusInternalServerError // сервер недоступен
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode
	}

	var resp_body ast.AstNode
	if err := json.NewDecoder(resp.Body).Decode(&resp_body); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, http.StatusInternalServerError // ошибка декодирования
	}

	return &resp_body, resp.StatusCode
}

func sendResult(taskID int, result float64, err string) {
	data := &Result{ID: taskID, Result: result, Error: err}
	jsonData, _ := json.Marshal(data)

	resp, er := http.Post(
		"http://localhost:8080/internal/task",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if er != nil {
		log.Printf("server error: %s", err)
	} else {
		resp.Body.Close()
	}
}

func worker(cfg config.Config) {
	for {
		task, code := getTask()
		if code == http.StatusNotFound {
			time.Sleep(2 * time.Second)
			continue
		}

		if code != http.StatusOK && code != http.StatusInternalServerError {
			log.Printf("error with code: %v", code)
			continue
		}

		if task == nil || task.Left == nil || task.Right == nil {
			continue
		}

		log.Printf("worker got expression with id %v", task.ID)
		result, err := calculate(task.Left.Value, task.Right.Value, task.Value, cfg)
		sendResult(task.ID, result, err)
		log.Printf("worker sent result %v with id %v", result, task.ID)
	}
}

func calculate(a, b string, operator string, cfg config.Config) (float64, string) {
	a_float, _ := strconv.ParseFloat(a, 64)
	b_float, _ := strconv.ParseFloat(b, 64)

	switch operator {
	case "*":
		time.Sleep(cfg.TimeMultiplication)
		return a_float * b_float, ""
	case "/":
		time.Sleep(cfg.TimeDivision)
		if b_float == 0 {
			return 0, "division by zero"
		}
		return a_float / b_float, ""
	case "+":
		time.Sleep(cfg.TimeAddition)
		return a_float + b_float, ""
	case "-":
		time.Sleep(cfg.TimeSubtraction)
		return a_float - b_float, ""
	default:
		return 0, ""
	}
}
