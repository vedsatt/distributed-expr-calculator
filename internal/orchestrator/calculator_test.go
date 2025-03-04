package orchestrator

import (
	"testing"

	"github.com/vedsatt/calc_prl/pkg/ast"
)

func TestSendTasks(t *testing.T) {
	// создаем аст для теста
	root := &ast.AstNode{
		ID:      1,
		AstType: "+",
		Left:    &ast.AstNode{ID: 2, AstType: "number", Value: "2"},
		Right:   &ast.AstNode{ID: 3, AstType: "number", Value: "3"},
	}

	// инициализируем каналы и мапу
	tasks := make(chan *ast.AstNode, 1)
	currTasks := make(map[int]*ast.AstNode)
	currTasks[root.ID] = root
	currTasks[root.Left.ID] = root.Left
	currTasks[root.Right.ID] = root.Right

	// запускаем функцию
	sendTasks(root, tasks, currTasks)

	// проверяем, что задача была отправлена в канал
	select {
	case task := <-tasks:
		if task.AstType != "+" {
			t.Errorf("expected task with type +, got %s", task.AstType)
		}
	default:
		t.Error("no task was sent to the channel")
	}
}

func TestDeleteAndUpdate(t *testing.T) {
	// инициализируем глобальные переменные
	currTasks = make(map[int]*ast.AstNode)
	last_result = make(chan float64, 1)

	// создаем дерево
	root := &ast.AstNode{
		ID:      1,
		AstType: "+",
		Left:    &ast.AstNode{ID: 2, AstType: "number", Value: "2"},
		Right:   &ast.AstNode{ID: 3, AstType: "number", Value: "3"},
	}

	// заполняем мапу текущих задач
	currTasks[root.ID] = root
	currTasks[root.Left.ID] = root.Left
	currTasks[root.Right.ID] = root.Right

	// создаем результат для обновления
	res := Result{ID: 1, Result: 5}

	// запускаем функцию
	result := deleteAndUpdate(res)

	// проверяем, что результат корректный
	if result != 5 {
		t.Errorf("expected result 5, got %f", result)
	}

	// проверяем, что листья удалены из мапы
	if _, exists := currTasks[2]; exists {
		t.Error("left node was not deleted from currTasks")
	}
	if _, exists := currTasks[3]; exists {
		t.Error("right node was not deleted from currTasks")
	}

	// проверяем, что нода обновлена
	if root.AstType != "number" || root.Value != "5.000000" {
		t.Errorf("expected node to be updated to number with value 5.000000, got %s with value %s", root.AstType, root.Value)
	}
}

func TestFillMap(t *testing.T) {
	// инициализируем глобальные переменные
	currTasks = make(map[int]*ast.AstNode)

	// создаем аст для теста
	root := &ast.AstNode{
		ID:      1,
		AstType: "+",
		Left:    &ast.AstNode{ID: 2, AstType: "number", Value: "2"},
		Right:   &ast.AstNode{ID: 3, AstType: "number", Value: "3"},
	}

	// запускаем функцию
	fillMap(root)

	// проверяем, что мапа заполнена корректно
	if len(currTasks) != 3 {
		t.Errorf("expected 3 nodes in currTasks, got %d", len(currTasks))
	}

	if _, exists := currTasks[1]; !exists {
		t.Error("root node was not added to currTasks")
	}
	if _, exists := currTasks[2]; !exists {
		t.Error("left node was not added to currTasks")
	}
	if _, exists := currTasks[3]; !exists {
		t.Error("right node was not added to currTasks")
	}
}
