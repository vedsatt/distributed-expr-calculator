package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/vedsatt/calc_prl/pkg/ast"
	"github.com/vedsatt/calc_prl/pkg/database"
)

const port = ":8080"

type Orchestrator struct{}

func New() *Orchestrator {
	return &Orchestrator{}
}

var (
	base   = database.New()
	mu     sync.Mutex // Мьютекс для синхронизации доступа к результатам
	astKey = contextKey{"ast"}
)

type contextKey struct {
	name string
}

type Ast struct {
	ID  int
	Ast *ast.AstNode
}

type Request struct {
	Expression string `json:"expression"`
}

type Error struct {
	Res string `json:"error"`
}

type RespID struct {
	Id int `json:"id"`
}

func errorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.WriteHeader(statusCode)
	e := Error{Res: err}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func checkId(id string) bool {
	pattern := "^[0-9]+$"
	r := regexp.MustCompile(pattern)
	return r.MatchString(id)
}

func (o *Orchestrator) Run() {
	mux := http.NewServeMux()

	expr := http.HandlerFunc(ExpressionHandler)
	GetData := http.HandlerFunc(GetDataHandler)
	tasks := http.HandlerFunc(TaskHandler)

	// хендлеры для пользователя
	mux.Handle("/api/v1/calculate", loggingMiddleware(databaseMiddleware(expr)))
	mux.Handle("/api/v1/expressions/", loggingMiddleware(GetData))

	// хендлер для общения с агентом
	mux.Handle("/internal/task", tasks)

	log.Printf("Starting sevrer on port %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
