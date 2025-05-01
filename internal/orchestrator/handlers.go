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

func logsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, URL: %s", r.Method, r.URL)

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("Method: %s, completion time: %v", r.Method, duration)
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
		expression := NewExpression(astRoot)

		// добавляем выражение в базу данных и получаем id
		expKey := base.PostData()
		w.WriteHeader(http.StatusCreated)
		log.Printf("Adding expression to database")

		w.Header().Set("Content-Type", "application/json")
		id := RespID{Id: expKey}
		json.NewEncoder(w).Encode(id)

		// запускаем вычисления в фоновом режиме
		go func() {
			// контекст нужен, чтобы передать выражение следующему хендлеру
			ctx := context.Background()
			exp := &Expr{ID: expKey, Expr: expression}
			ctx = context.WithValue(ctx, exprKey, exp)
			log.Printf("Adding AST: %v to ctx", exp.ID)
			req := r.WithContext(ctx)

			next.ServeHTTP(w, req)
		}()
	})
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	exp, ok := r.Context().Value(exprKey).(*Expr)
	if !ok || exp.Expr == nil {
		log.Printf("Error: unable to get AST from context")
		w.WriteHeader(http.StatusInternalServerError)
		base.UpdateData(exp.ID, 0, "error")
		return
	}

	res, err := exp.Expr.calc()
	if err != nil {
		log.Printf("Expression %v: zero division error detected", exp.ID)
		base.UpdateData(exp.ID, 0, "zero devision error")
		return
	}
	log.Printf("Expression %v calculated sucsessfully", exp.ID)

	base.UpdateData(exp.ID, res, "done")
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/:")
	if checkId(id) {
		id_int, err := strconv.Atoi(id)
		if err != nil {
			err_str := fmt.Sprintf("%s", err)
			errorResponse(w, err_str, http.StatusInternalServerError)
			log.Printf("Code: %v, Internal server error", http.StatusInternalServerError)
			return
		}

		data, err := base.GetExpression(id_int)
		if err != nil {
			errorResponse(w, fmt.Sprintf("%s", err), http.StatusNotFound)
			log.Printf("Code: %v, Expression was not found", http.StatusNotFound)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(data))
		if err != nil {
			errorResponse(w, "error with json data", http.StatusInternalServerError)
			log.Printf("Code: %v, Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	data, err := base.GetData()
	if err != nil {
		errorResponse(w, "empty base", http.StatusInternalServerError)
		log.Printf("Code: %v, Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // устанавливаем статус 200 OK
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(data))
	if err != nil {
		errorResponse(w, "error with json data", http.StatusInternalServerError)
		log.Printf("Code: %v, Internal server error", http.StatusInternalServerError)
		return
	}
}
