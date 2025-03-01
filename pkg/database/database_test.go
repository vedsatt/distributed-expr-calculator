package database

import (
	"encoding/json"
	"testing"
)

func TestPostData(t *testing.T) {
	base := New()
	id := base.PostData()

	if _, exists := database[id]; !exists {
		t.Errorf("Expected data with id %d to be posted, but it was not found in the database", id)
	}

	if database[id].Status != "in process" {
		t.Errorf("Expected status to be 'in process', but got '%s'", database[id].Status)
	}

	if database[id].Result != "" {
		t.Errorf("Expected result to be empty, but got '%s'", database[id].Result)
	}
}

func TestGetData(t *testing.T) {
	base := New()
	base.PostData()

	data, err := base.GetData()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	var expressions [][]byte
	err = json.Unmarshal(data, &expressions)
	if err != nil {
		t.Fatalf("Expected no error unmarshalling data, but got %v", err)
	}

	if len(expressions) == 0 {
		t.Error("Expected non-empty expressions, but got empty")
	}
}

func TestGetDataEmptyDatabase(t *testing.T) {
	base := New()

	_, err := base.GetData()
	if err == nil {
		t.Error("Expected error for empty database, but got nil")
	}

	if err.Error() != "empty database" {
		t.Errorf("Expected error 'empty database', but got '%v'", err)
	}
}

func TestGetExpression(t *testing.T) {
	base := New()
	id := base.PostData()

	data, err := base.GetExpression(id)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	var exp Data
	err = json.Unmarshal(data, &exp)
	if err != nil {
		t.Fatalf("Expected no error unmarshalling data, but got %v", err)
	}

	if exp.Id != id {
		t.Errorf("Expected id %d, but got %d", id, exp.Id)
	}
}

func TestGetExpressionNotFound(t *testing.T) {
	base := New()

	_, err := base.GetExpression(999)
	if err == nil {
		t.Error("Expected error for non-existent expression, but got nil")
	}

	expectedErr := "no expression with 999 id"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%v'", expectedErr, err)
	}
}

func TestUpdateData(t *testing.T) {
	base := New()
	id := base.PostData()

	err := base.UpdateData(id, 42.0, "done")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if database[id].Result != "42.000000" {
		t.Errorf("Expected result to be '42.000000', but got '%s'", database[id].Result)
	}

	if database[id].Status != "done" {
		t.Errorf("Expected status to be 'done', but got '%s'", database[id].Status)
	}
}

func TestUpdateDataNotFound(t *testing.T) {
	base := New()

	err := base.UpdateData(999, 42.0, "done")
	if err == nil {
		t.Error("Expected error for non-existent expression, but got nil")
	}

	expectedErr := "no expression with 999 id"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%v'", expectedErr, err)
	}
}
