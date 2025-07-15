package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikhathfirdose1/vaultify/internal/api"
	"github.com/nikhathfirdose1/vaultify/internal/config"
	"github.com/nikhathfirdose1/vaultify/internal/storage"
)

func main() {
	fmt.Println("Vaultify service starting...")

	cfg, err := config.LoadConfig("config/vaultify.yml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	fmt.Printf("Loaded config. Server will run on port %d\n", cfg.Server.Port)

	// BEFORE using Encrypt
	err = storage.LoadOrCreateKey(cfg.Encryption.KeyPath)
	if err != nil {
		log.Fatal("Key error:", err)
	}

	//  Initialize in-memory store
	storage.InitStore()
	storage.SetWhitelist(cfg.Auth.Tokens)

	// // Test secret encryption + store
	// secret := []byte("hello-secure-world")
	// encryptedBlob, err := storage.Encrypt(secret)
	// if err != nil {
	// 	log.Fatal("Encryption failed:", err)
	// }
	// fmt.Println("Encrypted blob:", encryptedBlob)

	// storage.StoreSecret("db_pass", encryptedBlob, 60)

	// time.Sleep(2 * time.Second)

	// fetched, err := storage.FetchSecret("db_pass")
	// if err != nil {
	// 	log.Fatal("Fetch failed:", err)
	// }
	// fmt.Println("Fetched blob:", fetched)

	// decrypted, err := storage.Decrypt(fetched)
	// if err != nil {
	// 	log.Fatal("Decryption failed:", err)
	// }

	// fmt.Println("Decrypted:", string(decrypted))

	r := mux.NewRouter()
	r.HandleFunc("/store", api.StoreHandler).Methods("POST")
	r.HandleFunc("/fetch/{name}", api.FetchHandler).Methods("GET")

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Println("Serving on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
