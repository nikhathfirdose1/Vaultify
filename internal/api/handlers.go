package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nikhathfirdose1/vaultify/internal/db"
	"github.com/nikhathfirdose1/vaultify/internal/storage"
)

type StoreRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

// POST /store
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !validateToken(token) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("Handler started")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var req StoreRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	encrypted, err := storage.Encrypt([]byte(req.Value))
	if err != nil {
		http.Error(w, "encryption failed", http.StatusInternalServerError)
		return
	}
	db.StoreSecret(req.Name, encrypted, req.TTL)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "stored %s", req.Name)
}

// GET /fetch/{name}
func FetchHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !validateToken(token) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	name := mux.Vars(r)["name"]
	blob, err := db.FetchSecret(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	plain, err := storage.Decrypt(blob)
	if err != nil {
		http.Error(w, "decryption failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(plain)
}

func validateToken(header string) bool {
	t := strings.TrimPrefix(header, "Bearer ")
	for _, token := range storage.WhitelistTokens() {
		if token == t {
			return true
		}
	}
	return false
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
