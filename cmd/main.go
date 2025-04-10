package main

import (
    "log"
    "os"

    "github.com/dibsnvas/golang-2025/internal/delivery"
    "github.com/dibsnvas/golang-2025/internal/repository"
)

func main() {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=postgres dbname=sales_ops port=5432 sslmode=disable"
    }

    db, err := repository.NewDB(dsn)
    if err != nil {
        log.Fatalf("Failed to connect DB: %v", err)
    }

    r := delivery.SetupRouter(db)

    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
