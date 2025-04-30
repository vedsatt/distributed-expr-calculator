package orchestrator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/vedsatt/calc_prl/internal/models"
)

var (
	managerCh chan channels
	tasksCh   chan *models.AstNode
	resultsCh chan models.Result
)

type channels struct {
	tasks   *chan *models.AstNode
	results *chan models.Result
}

type expression struct {
	node      *models.AstNode
	tasks     chan *models.AstNode
	results   chan models.Result
	currTasks map[int]*models.AstNode
}

func StartManager() {
	tasksCh = make(chan *models.AstNode)
	resultsCh = make(chan models.Result)

	managerCh = make(chan channels)
	go channelsManager(managerCh)
}

// channelsManager получает пару из двух каналов, которые должны общаться между собой
// каждый раз, когда приходит таска из канала с задачами, менеджер создает элемент в мапе,
// где ключ - айди таски, а значение - структура со связанной друг с другом парой каналов
// после этого он отправляет таску в общий канал со всеми тасками, чтобы хендлер взял таску и отправил ее агенту
// каждый раз, когда приходит результат таски от агента, хендлер отправляет ее в канал менеджера
// менеджер берет результат из канала и ищет элемент мапы с каналами по айди результата
// если такой элемент есть - отправляет результат с канал результатов, связанный с каналом задач
func channelsManager(chans chan channels) {
	chanMap := make(map[int]channels)
	mu := &sync.Mutex{}

	for ch := range chans {
		go func(ch channels, chanMap map[int]channels, mu *sync.Mutex) {
			for {
				task, ok := <-*ch.tasks
				if !ok {
					return
				}

				mu.Lock()
				chanMap[task.ID] = ch
				mu.Unlock()
				tasksCh <- task
			}
		}(ch, chanMap, mu)

		go func(ch channels, chanMap map[int]channels, mu *sync.Mutex) {
			for {
				result, ok := <-resultsCh
				if !ok {
					return
				}
				if chans, ok := chanMap[result.ID]; ok {
					*chans.results <- result

					mu.Lock()
					delete(chanMap, result.ID)
					mu.Unlock()
				}
			}
		}(ch, chanMap, mu)
	}
}

func NewExpression(node *models.AstNode) *expression {
	tasks := make(chan *models.AstNode)
	results := make(chan models.Result)

	chans := channels{
		tasks:   &tasks,
		results: &results,
	}
	managerCh <- chans

	return &expression{
		node:      node,
		tasks:     tasks,
		results:   results,
		currTasks: make(map[int]*models.AstNode),
	}
}

func (e *expression) calc() (float64, error) {
	e.fillMap(e.node)
	var result float64
	for {
		// проходимся по дереву и находим ноды, у которых оба листка - числа
		sendTasks(e.node, e.tasks, e.currTasks)

		select {
		case res := <-e.results:
			if res.Error != "" {
				// очищаем каналы
				close(e.tasks)
				close(e.results)

				log.Printf("id: %v, res: %v, err: %v", res.ID, res.Result, res.Error)
				return 0, errors.New(res.Error)
			}

			result = e.deleteAndUpdate(res)
			log.Println("Updated tree with new result")
		default:
			if len(e.currTasks) == 0 {
				return result, nil
			}
			time.Sleep(10 * time.Millisecond) // избегаем busy loop
		}

		// если все задачи удалены - результат получен, а значит можно завершать функцию
		if len(e.currTasks) == 0 {
			close(e.tasks)
			close(e.results)
			return result, nil
		}
	}
}

func sendTasks(node *models.AstNode, tasks chan<- *models.AstNode, currTasks map[int]*models.AstNode) {
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

func (e *expression) fillMap(node *models.AstNode) {
	if node == nil {
		return
	}

	// заполняем мапу, где ключ - айди ноды, а значение - сама нода
	mu.Lock()
	e.currTasks[node.ID] = node
	mu.Unlock()

	// обходим дерево методом пост-ордера
	e.fillMap(node.Left)
	e.fillMap(node.Right)
}

func (e *expression) deleteAndUpdate(res models.Result) float64 {
	// когда мы получаем результат ноды, мы удаляем ее листья, а потом меняем ноду на число для дальнейших вычислений
	// так как мапа ссылается на ноду, то, взаимодействуя с элементом мапы, мы напрямую взаимодействуем с нодой

	// проверяем, можно ли обращаться к листьям ноды для их удаления
	node, exists := e.currTasks[res.ID]
	if !exists || node.Left == nil || node.Right == nil {
		return 0
	}

	left := node.Left.ID
	right := node.Right.ID

	delete(e.currTasks, left)
	delete(e.currTasks, right)

	node.Value = fmt.Sprintf("%f", res.Result)
	node.AstType = "number"
	node.Left = nil
	node.Right = nil
	log.Printf("Updated node with id %d", node.ID)

	// простая обработка финального значения выражения
	if len(e.currTasks) == 1 {
		result, _ := strconv.ParseFloat(node.Value, 64)
		delete(e.currTasks, res.ID)
		return result
	}

	return 0
}
