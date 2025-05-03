package database

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"testing"

	"github.com/vedsatt/calc_prl/internal/models"
)

func TestUsers(t *testing.T) {
	db := dbTestConnect(t)
	defer db.Store.Close()

	users := []*models.User{
		{
			Login:    "user1",
			Password: "qwerty",
			ID:       1,
		},
		{
			Login:    "user2",
			Password: "ytrewq",
			ID:       2,
		},
	}

	ctx := context.TODO()
	for _, u := range users {
		id, err := db.InsertUser(ctx, u)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if id != u.ID {
			t.Fatalf("expected user ID %d, but got %d", u.ID, int(id))
		}
	}

	// добавление уже существующего юзера
	u := &models.User{
		Login:    "user1",
		Password: "123",
	}
	_, err := db.InsertUser(ctx, u)
	if err == nil {
		t.Fatal("expected error, because user already exists")
	}

	// поиск пользователей по айди
	for _, tc := range users {
		user, err := db.SelectUserByLogin(ctx, tc.Login)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != tc.ID || user.Login != tc.Login || user.Password != tc.Password {
			t.Fatalf("incorrect user: expected %v, but got %v", tc, user)
		}
	}

	// поиск несуществующего юзера
	_, err = db.SelectUserByLogin(ctx, "non-existent user")
	if err == nil {
		t.Fatal("expected error: non-existent user, but got nothing")
	}
}

func TestExpressions(t *testing.T) {
	db := dbTestConnect(t)
	db.InsertUser(context.TODO(), &models.User{ID: 1, Login: "1"})
	db.InsertUser(context.TODO(), &models.User{ID: 2, Login: "2"})

	type testcase struct {
		expr   *models.Expression
		userID int
	}
	testcases := []testcase{
		{
			expr: &models.Expression{
				ID:         1,
				Expression: "2+2",
				Status:     "done",
				Result:     4,
			},
			userID: 1,
		},
		{
			expr: &models.Expression{
				ID:         2,
				Expression: "2-2",
				Status:     "in process",
				Result:     0.0,
			},
			userID: 1,
		},
		{
			expr: &models.Expression{
				ID:         3,
				Expression: "2-2",
				Status:     "in process",
				Result:     0.0,
			},
			userID: 2,
		},
	}

	ctx := context.TODO()
	for _, tc := range testcases {
		id, err := db.InsertExpression(ctx, tc.expr, tc.userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if int(id) != tc.expr.ID {
			t.Fatalf("wrong expression id, expected: %v, but got %v", tc.expr.ID, id)
		}
	}

	_, err := db.InsertExpression(ctx, &models.Expression{}, 3)
	if err == nil {
		t.Fatal("expected error: user does not exist")
	}

	expected := testcases[1]
	got, err := db.SelectExprByID(ctx, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Expression != expected.expr.Expression ||
		got.ID != expected.expr.ID ||
		got.Result != expected.expr.Result ||
		got.Status != expected.expr.Status {
		t.Fatalf("test failure: expected %v, but got %v", expected.expr, got)
	}

	_, err = db.SelectExprByID(ctx, 5)
	if err == nil {
		t.Fatal("expected error: expression with ID 5 does not exist, but got nothing")
	}

	err = db.UpdateExpression(ctx, 2, "done", 0.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expr, err := db.SelectExprByID(ctx, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expr.Result != 0.0 || expr.Status != "done" ||
		expr.Expression != testcases[1].expr.Expression {
		t.Fatal("wrong expression after update")
	}

	err = db.UpdateExpression(ctx, 5, "done", 0.0)
	if err == nil {
		t.Fatal("expected error: expression does not exist, but got nothing")
	}
}

func dbTestConnect(t *testing.T) *SqlDB {
	t.Helper()
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	err = createTables(ctx, db)
	if err != nil {
		t.Fatalf("Failed tpo create tables: %v", err)
	}

	return &SqlDB{
		Store: db,
		usMu:  sync.Mutex{},
		expMu: sync.Mutex{},
	}
}
