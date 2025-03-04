package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vedsatt/calc_prl/internal/config"
	"github.com/vedsatt/calc_prl/pkg/ast"
)

// определяем пакетную переменную для базового URL
var baseURL = "http://localhost:8080"

// изменяем getTask и sendResult для использования baseURL
func get() (*ast.AstNode, bool) {
	resp, err := http.Get(baseURL + "/internal/task")
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	var resp_body ast.AstNode
	if err := json.NewDecoder(resp.Body).Decode(&resp_body); err != nil {
		return nil, false
	}

	return &resp_body, true
}

func sendRes(taskID int, result float64) {
	data := &Result{ID: taskID, Result: result}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(
		baseURL+"/internal/task",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// тесты
func TestGetTask(t *testing.T) {
	// создаем мок сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/task" {
			t.Errorf("Expected to request '/internal/task', got: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := &ast.AstNode{
			ID:    1,
			Value: "+",
			Left:  &ast.AstNode{Value: "2"},
			Right: &ast.AstNode{Value: "3"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// временно изменяем baseURL на адрес мок сервера
	oldBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = oldBaseURL }() // восстанавливаем значение после теста

	resp, ok := get()
	if !ok {
		t.Error("Expected to get a task, but got none")
	}

	if resp.ID != 1 || resp.Value != "+" || resp.Left.Value != "2" || resp.Right.Value != "3" {
		t.Errorf("Expected task with ID 1 and values 2 + 3, got: %v", resp)
	}
}

func TestSendResult(t *testing.T) {
	// создаем мок сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/task" {
			t.Errorf("Expected to request '/internal/task', got: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var result Result
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		if result.ID != 1 || result.Result != 5.0 {
			t.Errorf("Expected result with ID 1 and result 5.0, got: %v", result)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// временно изменяем baseURL на адрес mock-сервера
	oldBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = oldBaseURL }() // восстанавливаем значение после теста

	sendRes(1, 5.0)
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		a            string
		b            string
		operator     string
		expected     float64
		expected_err string
	}{
		{"2", "3", "+", 5.0, ""},
		{"5", "2", "-", 3.0, ""},
		{"4", "3", "*", 12.0, ""},
		{"10", "2", "/", 5.0, ""},
		{"10", "0", "/", 0.0, "division by zero"}, // деление на ноль
	}

	for _, tt := range tests {
		result, err := calculate(tt.a, tt.b, tt.operator, config.Config{})
		if result != tt.expected || err != tt.expected_err {
			t.Errorf("calculate(%s, %s, %s) = %v; expected %v, expected err: %s", tt.a, tt.b, tt.operator, result, tt.expected, err)
		}
	}
}

func TestWorker(t *testing.T) {
	// Создаем mock-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/task" {
			t.Errorf("Expected to request '/internal/task', got: %s", r.URL.Path)
		}
		if r.Method == http.MethodGet {
			response := &ast.AstNode{
				ID:    1,
				Value: "+",
				Left:  &ast.AstNode{Value: "2"},
				Right: &ast.AstNode{Value: "3"},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else if r.Method == http.MethodPost {
			var result Result
			if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
				t.Errorf("Error decoding request body: %v", err)
			}

			if result.ID != 1 || result.Result != 5.0 {
				t.Errorf("Expected result with ID 1 and result 5.0, got: %v", result)
			}

			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// временно изменяем baseURL на адрес мок сервера
	oldBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = oldBaseURL }() // восстанавливаем значение после теста

	go worker(config.Config{})
	time.Sleep(5 * time.Second) // Даем время воркеру выполнить свою работу
}
