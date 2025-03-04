package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vedsatt/calc_prl/pkg/ast"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, URL: %s", r.Method, r.URL)

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("Method: %s, completion time: %s", r.Method, duration)
	})
}

func databaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorResponse(w, "invalid request method", http.StatusMethodNotAllowed)
			log.Printf("Code: %v, Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var body Request
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorResponse(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Code: %v, Internal server error", http.StatusInternalServerError)
			return
		}

		// создаем аст по выражению заранее, чтобы обработать ошибки в мидлвеере
		astRoot, err := ast.Build(body.Expression)
		if err != nil {
			err_str := fmt.Sprintf("%s", err)
			errorResponse(w, err_str, http.StatusUnprocessableEntity)
			log.Printf("Code: %v, Unprocessable entity", http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("adding expression to database")
		expKey := base.PostData()

		w.Header().Set("Content-Type", "application/json")
		id := RespID{Id: expKey}
		json.NewEncoder(w).Encode(id)

		// с помощью контекста закидываем выражение в некст хендлер
		ast := &Ast{ID: expKey, Ast: astRoot}
		ctx := context.WithValue(r.Context(), astKey, ast)
		log.Printf("Adding AST: %v to ctx", *ast.Ast)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, URL: %s", r.Method, r.URL)

	ast, ok := r.Context().Value(astKey).(*Ast)
	if !ok || ast.Ast == nil {
		log.Printf("Error: unable to get AST from context")
		w.WriteHeader(http.StatusInternalServerError)
		base.UpdateData(ast.ID, 0, "error")
		return
	}

	fillMap(ast.Ast)
	log.Println("filled map")

	err := calc(ast.Ast)
	log.Println("end calculating")
	if err != "" {
		// сообщаем, что обнаружено деление на ноль
		log.Printf("Expression id: %v, zero division error detected", ast.ID)
		base.UpdateData(ast.ID, 0, "zero devision error")

		// очищаем каналы
		for len(tasks) > 0 {
			<-tasks
		}
		for len(results) > 0 {
			<-results
		}
		return
	}

	result := <-last_result

	// очищаем каналы
	for len(tasks) > 0 {
		<-tasks
	}
	for len(results) > 0 {
		<-results
	}
	base.UpdateData(ast.ID, result, "done")
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/:")
	if id != "" {
		id_int, err := strconv.Atoi(id)
		if err != nil {
			err_str := fmt.Sprintf("%s", err)
			errorResponse(w, err_str, http.StatusUnprocessableEntity)
			log.Printf("Code: %v, Unprocessable entity", http.StatusUnprocessableEntity)
			return
		}

		data, err := base.GetExpression(id_int)
		if err != nil {
			errorResponse(w, fmt.Sprintf("%s", err), http.StatusUnprocessableEntity)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(data))
		if err != nil {
			errorResponse(w, "error with json data", http.StatusInternalServerError)
			log.Printf("Code: %v, Unprocessable entity", http.StatusInternalServerError)
			return
		}
		return
	}

	data, err := base.GetData()
	if err != nil {
		errorResponse(w, "empty base", http.StatusUnprocessableEntity)
		log.Printf("Code: %v, Unprocessable entity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // устанавливаем статус 200 OK
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(data))
	if err != nil {
		errorResponse(w, "error with json data", http.StatusInternalServerError)
		log.Printf("Code: %v, Unprocessable entity", http.StatusInternalServerError)
		return
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		select {
		case task := <-tasks:
			// если есть задача, отправляем её агенту
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		default:
			// если задач нет, возвращаем http 404
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPost:
		var result Result
		json.NewDecoder(r.Body).Decode(&result)
		results <- result

		// закрываем тело запроса
		defer r.Body.Close()
	}
}
