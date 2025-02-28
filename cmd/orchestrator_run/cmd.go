package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/vedsatt/calc_prl/pkg/ast"
)

type Result struct {
	ID     int
	Result float64
}

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

var (
	product time.Duration = 2 * time.Millisecond
	plus    time.Duration = 1 * time.Millisecond
)

func main() {
	expression := "1 * (2 - 3) + (4 - 5 + 3) + 1 - (4 * 4 + 20) - (2 * (3 - 5))"
	astRoot, err := ast.Build(expression)
	if err != nil {
		fmt.Println(err)
	}

	tasks := make(chan *ast.AstNode, 10)
	results := make(chan Result, 10)
	last_result := make(chan float64, 1)

	currTasks := make(map[int]*ast.AstNode)

	fillMap(currTasks, astRoot)

	wg.Add(2)

	// оркестратор
	go func() {
		defer wg.Done()
		calc(astRoot, tasks, results, currTasks, last_result)
	}()

	//агент
	go func() {
		defer wg.Done()
		agent(tasks, results)
	}()

	wg.Wait()
	close(results)

	result := <-last_result
	fmt.Println(result)
}

func calc(node *ast.AstNode, tasks chan *ast.AstNode, results chan Result, currTasks map[int]*ast.AstNode, last_result chan float64) {
	defer close(tasks)
	var result float64
	for {
		// проходимся по дереву и находим ноды, у которых оба листка - числа
		sendTasks(node, tasks, currTasks)

		select {
		case res, ok := <-results:
			if !ok {
				last_result <- result
				return // Канал results закрыт
			}
			result = delete_and_update(res, currTasks)
		default:
			if len(currTasks) == 0 {
				last_result <- result
				return
			}
			time.Sleep(10 * time.Millisecond) // Избегаем busy loop
		}

		// если все задачи удалены - результат получен, а значит можно завершать функцию
		if len(currTasks) == 0 {
			last_result <- result
			return
		}
	}
}

func sendTasks(node *ast.AstNode, tasks chan<- *ast.AstNode, currTasks map[int]*ast.AstNode) {
	if node == nil {
		return
	}

	if node.AstType == "number" {
		return
	}

	// проверяем, что узел не обработан, а его листья - числа
	if node.Left != nil && node.Right != nil &&
		node.Left.AstType == "number" && node.Right.AstType == "number" {
		if node, exists := currTasks[node.ID]; exists && !node.Counting {
			node.Counting = true
			tasks <- node
		}
	}

	// пост-ордер для отправки всех готовых тасков агенту
	sendTasks(node.Left, tasks, currTasks)
	sendTasks(node.Right, tasks, currTasks)
}

func fillMap(currTasks map[int]*ast.AstNode, node *ast.AstNode) {
	if node == nil {
		return
	}

	// заполняем мапу, где ключ - айди ноды, а значение - сама нода
	mu.Lock()
	currTasks[node.ID] = node
	mu.Unlock()

	// обходим дерево методом пост-ордера
	fillMap(currTasks, node.Left)
	fillMap(currTasks, node.Right)
}

func delete_and_update(res Result, currTasks map[int]*ast.AstNode) float64 {
	// когда мы получаем результат ноды, мы удаляем ее листья, а потом меняем ноду на число для дальнейших вычислений
	// так как мапа ссылается на ноду, то, взаимодействуя с элементом мапы, мы напрямую взаимодействуем с нодой
	mu.Lock()
	defer mu.Unlock()

	// проверяем, можно ли обращаться к листьям ноды для их удаления
	node, exists := currTasks[res.ID]
	if !exists || node.Left == nil || node.Right == nil {
		return 0
	}

	left := node.Left.ID
	right := node.Right.ID

	delete(currTasks, left)
	delete(currTasks, right)

	node.Value = fmt.Sprintf("%f", res.Result)
	node.AstType = "number"
	node.Left = nil
	node.Right = nil

	if len(currTasks) == 1 {
		result, _ := strconv.ParseFloat(node.Value, 64)
		delete(currTasks, res.ID)
		return result
	}

	return 0
}

func agent(tasks chan *ast.AstNode, results chan Result) {
	power := 2
	var wg sync.WaitGroup
	wg.Add(power)

	for i := 0; i < power; i++ {
		go func() {
			defer wg.Done()
			for task := range tasks {
				if task == nil || task.Left == nil || task.Right == nil {
					continue
				}
				result := count(task.Left.Value, task.Right.Value, task.Value)
				results <- Result{ID: task.ID, Result: result}
			}
		}()
	}

	wg.Wait()
}

// func worker(tasks chan *ast.AstNode, results chan Result) {
// 	for {
// 		select {
// 		case task := <-tasks:
// 			result := count(task.Left.Value, task.Right.Value, task.Value)
// 			res := Result{ID: task.ID, Result: result}
// 			results <- res
// 		default:
// 		}
// 	}
// }

func worker(tasks <-chan *ast.AstNode, results chan<- Result) {
	for task := range tasks { // завершится при закрытии tasks
		if task == nil || task.Left == nil || task.Right == nil {
			continue
		}

		result := count(task.Left.Value, task.Right.Value, task.Value)
		fmt.Println(Result{ID: task.ID, Result: result})
		results <- Result{ID: task.ID, Result: result}
	}
}

func count(a, b string, operator string) float64 {
	a_float, _ := strconv.ParseFloat(a, 64)
	b_float, _ := strconv.ParseFloat(b, 64)

	if operator == "*" || operator == "/" {
		time.Sleep(product)
	} else {
		time.Sleep(plus)
	}

	switch operator {
	case "*":
		return a_float * b_float
	case "/":
		return a_float / b_float
	case "+":
		return a_float + b_float
	case "-":
		return a_float - b_float
	default:
		return 0
	}
}
