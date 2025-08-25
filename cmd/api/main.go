package main

import (
	"log"

	"github.com/MohammadTaghipour/social/internal/env"
	"github.com/MohammadTaghipour/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil) // TODO: pass a real db connection

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
