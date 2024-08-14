package main

import (
	"net/http"
	"net/http/httptest"
    "testing"
	"github.com/gorilla/mux"
)

// Test happy path of submitting a well-formed GET /key/{length} request
func TestGetKeyHandlerValid(t *testing.T) {
    req, err := http.NewRequest("GET", "/key/100", nil)
	
	if err != nil {
        t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{length}", getRandomKey).Methods("GET")
	router.ServeHTTP(rr, req)

	// Checks for 200 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getRandomKey returned wrong status code: got %v want %v",
  			status, http.StatusOK)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v",
			ctype, "application/json")
	}
}

// Test unhappy path of an invalid length in GET /key/{length} request
func TestGetKeyHandlerInvalid(t *testing.T) {
    req, err := http.NewRequest("GET", "/key/sddafypqx", nil)
	
	if err != nil {
        t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{length}", getRandomKey).Methods("GET")
	router.ServeHTTP(rr, req)

	// Checks for 400 status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("getRandomKey returned wrong status code: got %v want %v",
  			status, http.StatusBadRequest)
	}
}

// Test unhappy path of an out-of-range length in GET /key/{length} request
func TestGetKeyHandlerOutOfRange(t *testing.T) {
    req, err := http.NewRequest("GET", "/key/1050", nil)
	
	if err != nil {
        t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{length}", getRandomKey).Methods("GET")
	router.ServeHTTP(rr, req)

	// Checks for 400 status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("getRandomKey returned wrong status code: got %v want %v",
  			status, http.StatusBadRequest)
	}
}