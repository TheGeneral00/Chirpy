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

	req, err := http.NewRequest("Get", "http://localhost/users", nil)
	if err != nil {
		t.Fatalf("Failed to create requet with error: %s", err)
	}
	uuid := uuid.New().String()
	req.Header.Set("X-User-ID", uuid)
	w := httptest.NewRecorder()

	//Add function to read time of the received request to guess time in the db
	rows := sqlmock.NewRows([]string{"id", "user_id", "method", "method_details", "created_at"}).AddRow(1, uuid, "GET", "http://localhost/users", )
	mock.ExpectQuery("^SELECT (.+) FROM user_events$").WillReturnRows(rows)

	//Add a function call or what ever works here to call the AddUser handler on the test db 
	handler := testCfg.MiddlewareCreateUserEvent(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(w, req)
   	
	mock.ExpectExec("INSERT INTO user_events").WithArgs(sqlmock.AnyArg(), "GET", "/user_events").WillReturnResult(sqlmock.NewResult(1,1))
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("DB assoziation failed: %s", err)
	}
}
