package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
)

func TestCreateUser(t *testing.T){
	body := helpers.
	httptest.NewRequest(http.MethodPost, "/api/users", body)
}
