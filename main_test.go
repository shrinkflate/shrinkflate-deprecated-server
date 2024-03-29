package main

import (
	"github.com/aerogo/aero"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticRoutes(t *testing.T) {
	app := configure(aero.New())
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Invalid status %d", response.Code)
	}
}
