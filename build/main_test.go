package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	// helper "bitbucket.org/michaelchandrag/kumparan-test/helper"
)

func TestHello(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	router.ServeHTTP(w, req)

	status := w.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	type response struct {
		Message 	string 		`json:"message"`
	}
	data := &response{}
	err := json.NewDecoder(w.Body).Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != "hello" {
		t.Fatal("Message not hello.")
	}
}