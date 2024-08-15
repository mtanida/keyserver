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
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
    NUM_BUCKETS int = 20

    max_key_size int = 1024

    // Collection of Prometheus counters for response status codes
    statusCodeCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "keyserver_status_codes_count",
            Help: "Count of status codes returned in GET /key requests.",
        },
        []string{"status_code"},
    )

    // Prometheus key len histogram
    keyLenHistogram prometheus.Histogram
)

// Initialize and register key len histogram
func initKeyHistogram() {
    // Create buckets
    bucketWidth := float64(max_key_size) / float64(NUM_BUCKETS)

    // Register Prometheus key length histogram metric
    keyLenHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
        Name: "keyserver_length_distribution",
        Help: "Histogram of key lengths in GET /key requests.",
        Buckets: prometheus.LinearBuckets(0, bucketWidth, NUM_BUCKETS),
    })
}

// Handler for GET /key request. Returns a JSON object with request
// bytes encoded in base64
func getRandomKey(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Get path parameter
    s := mux.Vars(r)["length"]
    length, err := strconv.Atoi(s)

    // Check for string to int conversion error
    if err != nil {
        log.Print(err)
        w.WriteHeader(http.StatusBadRequest)
        statusCodeCounter.WithLabelValues(fmt.Sprint(http.StatusBadRequest)).Inc()
        errorMsg := map[string]string{"error": "Invalid length received."}
        json.NewEncoder(w).Encode(errorMsg)
        return
    }

    // Ensure 0 < length < max_key_size
    if (length < 0 || length > max_key_size) {
        log.Printf("received out of range length of %d\n", length)
        w.WriteHeader(http.StatusBadRequest)
        statusCodeCounter.WithLabelValues(fmt.Sprint(http.StatusBadRequest)).Inc()
        errStr := fmt.Sprintf("length must be in range [0, %d].", max_key_size)
        errorMsg := map[string]string{"error": errStr}
        json.NewEncoder(w).Encode(errorMsg)
        return
    }

    // Generate random key
    // from https://www.tutorialspoint.com/how-to-generate-random-string-characters-in-golang
    bytes := make([]byte, length)
    _, keyErr := rand.Read(bytes)
    if keyErr != nil {
        log.Print(keyErr)
        w.WriteHeader(http.StatusInternalServerError)
        statusCodeCounter.WithLabelValues(fmt.Sprint(http.StatusInternalServerError)).Inc()
        errorMsg := map[string]string{"error": "Encountered error generating key."}
        json.NewEncoder(w).Encode(errorMsg)
        return
    }
    randKey := base64.StdEncoding.EncodeToString(bytes)

    // Return key
    w.WriteHeader(http.StatusOK)
    statusCodeCounter.WithLabelValues(fmt.Sprint(http.StatusOK)).Inc()
    keyLenHistogram.Observe(float64(length))
    returnMsg := map[string]string{"key": randKey}
    json.NewEncoder(w).Encode(returnMsg)
}

func main() {
    // Process command line args
    var port int
    flag.IntVar(&max_key_size, "max-size", 1024, "maximum key size")
    flag.IntVar(&port, "srv-port", 1123, "server listening port")
    flag.Parse()

    // Initialize prometheus histogram
    initKeyHistogram()

    // Set up route handlers
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/key/{length}", getRandomKey).Methods("GET")

    // Setup metrics endpoint
    router.Handle("/metrics", promhttp.Handler())

    // Run server
    fmt.Printf("Starting server at port %d...\n", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
