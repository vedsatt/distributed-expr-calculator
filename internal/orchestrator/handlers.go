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
			log.Printf("Code: %v, error with decoding request body", http.StatusInternalServerError)
			return
		}

		// добавляем выражение в базу данных и получаем id
		log.Printf("Adding expression to database")
		expID := base.PostData()
		respID := RespID{Id: expID}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respID)

		// запускаем вычисления в фоновом режиме
		go func() {
			expr := &ExprReq{
				exp: body.Expression,
				id:  expID,
			}
			ctx := context.WithValue(r.Context(), ctxKey, expr)
			next.ServeHTTP(w, r.WithContext(ctx))
		}()
	})
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	expr := r.Context().Value(ctxKey).(*ExprReq)

	// создаем аст по выражению
	astRoot, err := ast.Build(expr.exp)
	if err != nil {
		err_str := fmt.Sprintf("%s", err)

		log.Printf("Expression %v: %s detected", expr.id, err_str)
		base.UpdateData(expr.id, 0.0, err_str)
		return
	}
	exp := NewExpression(astRoot)

	res, err := exp.calc()
	if err != nil {
		log.Printf("Expression %v: zero division error detected", expr.id)
		base.UpdateData(expr.id, 0, "zero division error")
		return
	}
	log.Printf("Expression %v calculated sucsessfully", expr.id)

	base.UpdateData(expr.id, res, "done")
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
