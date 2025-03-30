package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Data struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

type Expressions struct {
	Expr []Data `json:"expressions"`
}

type Expression struct {
	Expr Data `json:"expression"`
}

type Base struct{}

func New() *Base {
	database = make(map[int]*Data)
	return &Base{}
}

var database map[int]*Data

// формируем новое выражение и получаем его айди
func (b *Base) PostData() int {
	id := int(time.Now().UnixNano())

	database[id] = &Data{Id: id, Status: "in process", Result: ""}

	return id
}

// получаем все выражения
func (b *Base) GetData() ([]byte, error) {
	if len(database) == 0 {
		return nil, errors.New("empty database")
	}

	var e Expressions
	for _, exp := range database {
		data := exp
		e.Expr = append(e.Expr, *data)
	}

	jsonData, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при сериализации в JSON:", err)
		return nil, err
	}
	return jsonData, nil
}

// получаем выражение по айди
func (b *Base) GetExpression(id int) ([]byte, error) {
	exp, ok := database[id]

	if !ok {
		err := fmt.Sprintf("no expression with id: %d", id)
		return nil, errors.New(err)
	}
	data := Expression{Expr: *exp}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// обновляем результат и статус выражения по айди
func (b *Base) UpdateData(id int, res float64, status string) error {
	if _, ok := database[id]; !ok {
		err := fmt.Sprintf("no expression with %d id", id)
		return errors.New(err)
	}

	res_string := fmt.Sprintf("%.2f", res)
	database[id].Result = res_string
	database[id].Status = status

	return nil
}
