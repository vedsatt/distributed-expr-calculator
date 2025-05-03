package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vedsatt/calc_prl/internal/models"
)

type (
	SqlDB struct {
		Store *sql.DB
		usMu  sync.Mutex
		expMu sync.Mutex
	}

	Expressions struct {
		Expressions []models.Expression `json:"expressions"`
	}
)

func createTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT UNIQUE,
		password TEXT NOT NULL
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		expression TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		status TEXT,
		result REAL,
	
		FOREIGN KEY (user_id)  REFERENCES users (id)
	);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}

func NewDB() *SqlDB {
	log.Println("Opening database connection...")
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "data/store.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	log.Println("Pinging database...")
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Creating tables...")
	err = createTables(ctx, db)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")
	return &SqlDB{
		Store: db,
		usMu:  sync.Mutex{},
		expMu: sync.Mutex{},
	}
}

// создаем нового пользователя
func (db *SqlDB) InsertUser(ctx context.Context, user *models.User) (int64, error) {
	var q = `
	INSERT INTO users (name, password) values ($1, $2)
	`
	db.usMu.Lock()
	result, err := db.Store.ExecContext(ctx, q, user.Name, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	db.usMu.Unlock()

	return id, nil
}

// формируем новое выражение и получаем его айди
func (db *SqlDB) InsertExpression(ctx context.Context, expression *models.Expression, userID int) (int64, error) {
	var q = `
	INSERT INTO expressions (expression, user_id, status, result) values ($1, $2, $3, $4)
	`

	db.expMu.Lock()
	result, err := db.Store.ExecContext(ctx, q, expression.Expression, userID, expression.Status, expression.Result)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	db.expMu.Unlock()

	return id, nil
}

// получаем все выражения конкретного пользователя, формат - json файл с массивом json выражений
func (db *SqlDB) SelectExpressions(ctx context.Context, userID int) ([]byte, error) {
	var expressions []models.Expression
	var q = "SELECT id, expression, status, result FROM expressions WHERE user_id = $1"

	rows, err := db.Store.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := models.Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.Status, &e.Result)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	jsonData, err := json.MarshalIndent(expressions, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// находим конкретного пользователя по айди
func (db *SqlDB) SelectUserByID(ctx context.Context, id int) (models.User, error) {
	u := models.User{}
	var q = "SELECT id, name, password FROM users WHERE id = $1"
	err := db.Store.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Name, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil
}

// находим конкретного пользователя по логину
func (db *SqlDB) SelectUserByLogin(ctx context.Context, login string) (models.User, error) {
	u := models.User{}
	var q = "SELECT id, name, password FROM users WHERE name = $1"
	err := db.Store.QueryRowContext(ctx, q, login).Scan(&u.ID, &u.Name, &u.Password)
	if err != nil {
		return u, err
	}
	return u, nil
}

// находим конкретное выражение
func (db *SqlDB) SelectExprByID(ctx context.Context, id int) (models.Expression, error) {
	e := models.Expression{}
	var q = "SELECT id, expression, status, result FROM expressions WHERE id = $1"
	err := db.Store.QueryRowContext(ctx, q, id).Scan(&e.ID, &e.Expression, &e.Status, &e.Result)
	if err != nil {
		return e, err
	}

	return e, nil
}

// обновляем результат и статус выражения по айди
func (db *SqlDB) UpdateExpression(ctx context.Context, id int, status string, result float64) error {
	var q = "UPDATE expressions SET status = $1, result = $2 WHERE id = $3"
	_, err := db.Store.ExecContext(ctx, q, status, result, id)
	if err != nil {
		return err
	}

	return nil
}
