package main

import (
	"github.com/MohammadTaghipour/social/internal/db"
	"github.com/MohammadTaghipour/social/internal/env"
	"github.com/MohammadTaghipour/social/internal/store"
)

// run using: go run cmd/migrate/seed/main.go
func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	maxOpenConns := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleConns := env.GetInt("DB_MAX_IDLE_CONNS", 30)
	maxIdleTime := env.GetString("DB_MAX_IDLE_TIME", "15m")

	conn, _ := db.New(
		addr,
		maxOpenConns,
		maxIdleConns,
		maxIdleTime,
	)

	store := store.NewStorage(conn)
	db.Seed(store)
}
