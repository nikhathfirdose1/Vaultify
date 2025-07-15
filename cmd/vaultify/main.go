package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nikhathfirdose1/vaultify/internal/api"
	"github.com/nikhathfirdose1/vaultify/internal/config"
	"github.com/nikhathfirdose1/vaultify/internal/db"
	"github.com/nikhathfirdose1/vaultify/internal/metrics"
	"github.com/nikhathfirdose1/vaultify/internal/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("Vaultify service starting...")

	cfg, err := config.LoadConfig("config/vaultify.yml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	fmt.Printf("Loaded config. Server will run on port %d\n", cfg.Server.Port)


	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
	cfg.Database.Password, cfg.Database.Name)

	if err := db.InitDB(connStr); err != nil {
	log.Fatal("DB connection failed:", err)
	}

	err = storage.LoadOrCreateKey(cfg.Encryption.KeyPath)
	if err != nil {
		log.Fatal("Key error:", err)
	}

	storage.InitStore()
	storage.SetWhitelist(cfg.Auth.Tokens)
	metrics.InitMetrics()

	r := mux.NewRouter()
	r.Use(accessLogger(cfg.Server.LogPath))
	r.HandleFunc("/store", api.StoreHandler).Methods("POST")
	r.HandleFunc("/fetch/{name}", api.FetchHandler).Methods("GET")
	r.HandleFunc("/healthz", api.HealthCheckHandler).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())


	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Println("Serving on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}



func accessLogger(logPath string) mux.MiddlewareFunc {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Could not open access log file: %v", err)
	}
	logger := log.New(file, "", log.LstdFlags)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			res := &responseWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(res, r)
			duration := time.Since(start)
			logger.Printf("%s %s %d %s\n", r.Method, r.URL.Path, res.status, duration)

			// Update Prometheus metric
			metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

