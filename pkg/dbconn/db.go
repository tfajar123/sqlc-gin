package dbconn

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() *sql.DB {
    connStr := "postgres://postgres:root@localhost:5432/sqlc_migrate?sslmode=disable"

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to open DB:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }

    fmt.Println("DB connection successful")

    return DB
}