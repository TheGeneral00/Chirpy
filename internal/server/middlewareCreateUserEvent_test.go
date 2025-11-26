package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
)

func TestCreateUserEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to spin up mock db: %s", err)
	}
	defer db.Close()
	
	queries := database.New(db)
	testCfg := &APIConfig{DBQueries: queries}

	mock.ExpectExec("INSERT INTO user_events").WithArgs(sqlmock.AnyArg(), "GET", "/test").WillReturnResult(sqlmock.NewResult(1,1))

	req := httptest.NewRequest("Get", "/test", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	w := httptest.NewRecorder()

	handler := testCfg.MiddlewareCreateUserEvent(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("DB assoziation failed: %s", err)
	}
}
