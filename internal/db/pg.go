package db

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) error {
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
    return DB.Ping()
}

func StoreSecret(name string, blob []byte, ttl int) error {
    _, err := DB.Exec(`
        INSERT INTO secrets (name, blob, ttl_seconds, created_at)
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (name)
        DO UPDATE SET blob = EXCLUDED.blob, ttl_seconds = EXCLUDED.ttl_seconds, created_at = NOW()
    `, name, blob, ttl)
    return err
}

func FetchSecret(name string) ([]byte, error) {
    var blob []byte
    var created time.Time
    var ttl int

    err := DB.QueryRow(`
        SELECT blob, created_at, ttl_seconds FROM secrets WHERE name = $1
    `, name).Scan(&blob, &created, &ttl)
    if err != nil {
        return nil, err
    }

    if time.Since(created) > time.Duration(ttl)*time.Second {
        return nil, fmt.Errorf("secret expired")
    }

    return blob, nil
}
