package orchestrator

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/vedsatt/calc_prl/pkg/ast"
)

type Result struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
	Error  string  `json:"error"`
}

var (
	tasks       chan *ast.AstNode
	results     chan Result
	last_result chan float64
	currTasks   map[int]*ast.AstNode
)

func init() {
	last_result = make(chan float64, 1)
	tasks = make(chan *ast.AstNode)
	results = make(chan Result)
	currTasks = make(map[int]*ast.AstNode)
}

func calc(node *ast.AstNode) string {
	var result float64
	for {
		// проходимся по дереву и находим ноды, у которых оба листка - числа
		sendTasks(node, tasks, currTasks)

		select {
		case res := <-results:
			if res.Error != "" {
				last_result <- 0
				currTasks = make(map[int]*ast.AstNode)
				return res.Error
			}
			log.Printf("id: %v, res: %v, err: %v", res.ID, res.Result, res.Error)
			result = delete_and_update(res)
			log.Println("Updated tree with new result")
		default:
			if len(currTasks) == 0 {
				last_result <- result
				return ""
			}
			time.Sleep(10 * time.Millisecond) // избегаем busy loop
		}

		// если все задачи удалены - результат получен, а значит можно завершать функцию
		if len(currTasks) == 0 {
			last_result <- result
			return ""
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

func fillMap(node *ast.AstNode) {
	if node == nil {
		return
	}

	// заполняем мапу, где ключ - айди ноды, а значение - сама нода
	mu.Lock()
	currTasks[node.ID] = node
	mu.Unlock()

	// обходим дерево методом пост-ордера
	fillMap(node.Left)
	fillMap(node.Right)
}

func delete_and_update(res Result) float64 {
	// когда мы получаем результат ноды, мы удаляем ее листья, а потом меняем ноду на число для дальнейших вычислений
	// так как мапа ссылается на ноду, то, взаимодействуя с элементом мапы, мы напрямую взаимодействуем с нодой

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
	log.Printf("Updated node with id %d", node.ID)

	// простая обработка финального значения выражения
	if len(currTasks) == 1 {
		result, _ := strconv.ParseFloat(node.Value, 64)
		delete(currTasks, res.ID)
		return result
	}

	return 0
}
