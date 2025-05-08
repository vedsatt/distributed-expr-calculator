package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/vedsatt/calc_prl/pkg/database"
)

const port = ":8080"

type (
	Orchestrator struct {
	}

	ExpressionReq struct {
		Expression string `json:"expression"`
	}

	RespID struct {
		Id int `json:"id"`
	}

	ErrorResponse struct {
		Res string `json:"error" example:"Internal server error"`
	}

	Expression struct {
		exp string
		id  int
	}

	contextKey string
	userid     string
)

func New() *Orchestrator {
	return &Orchestrator{}
}

var (
	db     *database.SqlDB
	mu     sync.Mutex // Мьютекс для синхронизации доступа к результатам
	ctxKey contextKey = "expression id"
	userID userid     = "user id"
)

func checkCookie(cookie *http.Cookie, err error) bool {
	if err != nil {
		return false
	}

	token := cookie.Value
	return !(len(token) == 0)
}

func errorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.WriteHeader(statusCode)
	e := ErrorResponse{Res: err}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func checkId(id string) bool {
	if id == "-1" {
		return false
	}

	pattern := "^[0-9]+$"
	r := regexp.MustCompile(pattern)
	return r.MatchString(id)
}

func (o *Orchestrator) Run() {
	// подключение к бд
	db = database.NewDB()
	defer db.Store.Close()

	// запуск менеджера каналов выражений
	StartManager()
	// запуск сервера для общения с агентом
	go runGRPC()

	mux := http.NewServeMux()

	register := http.HandlerFunc(RegisterHandler)
	login := http.HandlerFunc(LoginHandler)
	expr := http.HandlerFunc(ExpressionHandler)
	getData := http.HandlerFunc(GetDataHandler)

	// хендлеры
	mux.Handle("/api/v1/register", logsMiddleware(register))
	mux.Handle("/api/v1/login", logsMiddleware(login))
	mux.Handle("/api/v1/calculate", logsMiddleware(authMiddleware(databaseMiddleware(expr))))
	mux.Handle("/api/v1/expressions/", logsMiddleware(authMiddleware(getData)))

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, mux))

}
