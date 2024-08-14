package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

var max_key_size int = 1024


func getRandomKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get path parameter
	s := mux.Vars(r)["length"]
	length, err := strconv.Atoi(s)

	// Check for string to int conversion error
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := map[string]string{"error": "Invalid length received"}
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Ensure 0 < length < max_key_size
	if (length < 0 || length > max_key_size) {
		fmt.Printf("received out of range length of %d\n", length)
		w.WriteHeader(http.StatusBadRequest)
		errStr := fmt.Sprintf("length must be in range [0, %d]", max_key_size)
		errorMsg := map[string]string{"error": errStr}
		json.NewEncoder(w).Encode(errorMsg)
	    return   
	}

    // Generate random key
	// from https://www.tutorialspoint.com/how-to-generate-random-string-characters-in-golang
    bytes := make([]byte, length)
	_, keyErr := rand.Read(bytes)
	if keyErr != nil {
		fmt.Println(keyErr)
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := map[string]string{"error": "Encountered error generating key"}
		json.NewEncoder(w).Encode(errorMsg)
	    return   
	}
	randKey := base64.StdEncoding.EncodeToString(bytes)

	// Return key
    w.WriteHeader(http.StatusOK)
    returnMsg := map[string]string{"key": randKey}
	json.NewEncoder(w).Encode(returnMsg)
}

func main() {
	// Process command line args
	var port int
    flag.IntVar(&max_key_size, "max-size", 1024, "maximum key size")
    flag.IntVar(&port, "srv-port", 1123, "server listening port")
	flag.Parse()

	// Set up route handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{length}", getRandomKey).Methods("GET")

    // Run server
	fmt.Printf("Starting server at port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}