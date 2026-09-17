package main

import (
	"bytes"
	"database/sql"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func setupCreateOrderRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/orders", CreateOrder(db))

	return router
}

func TestCreateOrderMissingID(t *testing.T) {
	router := setupCreateOrderRouter(nil)

	body := bytes.NewBufferString(
		`{"user_id":"test-user","status":"pending"}`,
	)

	req := httptest.NewRequest("POST", "/orders", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateOrderInvalidStatus(t *testing.T) {
	router := setupCreateOrderRouter(nil)

	body := bytes.NewBufferString(
		`{"id":"order-123","user_id":"test-user","status":"banana"}`,
	)

	req := httptest.NewRequest("POST", "/orders", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateOrderInvalidJSON(t *testing.T) {
	router := setupCreateOrderRouter(nil)

	body := bytes.NewBufferString(`{"id":`)

	req := httptest.NewRequest("POST", "/orders", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateOrderSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO orders (id, user_id, status) VALUES ($1, $2, $3)",
	)).
		WithArgs("order-123", "test-user", "pending").
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := setupCreateOrderRouter(db)

	body := bytes.NewBufferString(
		`{"id":"order-123","user_id":"test-user","status":"pending"}`,
	)

	req := httptest.NewRequest("POST", "/orders", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected 201, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}

func TestCreateOrderDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO orders (id, user_id, status) VALUES ($1, $2, $3)",
	)).
		WithArgs("order-123", "test-user", "pending").
		WillReturnError(sql.ErrConnDone)

	router := setupCreateOrderRouter(db)

	body := bytes.NewBufferString(
		`{"id":"order-123","user_id":"test-user","status":"pending"}`,
	)

	req := httptest.NewRequest("POST", "/orders", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected 500, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}
