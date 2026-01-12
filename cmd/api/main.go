package main

import (
	"fmt"
	"log"

	"github.com/runAlgo/notes-api/internal/config"
	"github.com/runAlgo/notes-api/internal/db"
	"github.com/runAlgo/notes-api/internal/server"
)

// config -> db -> router -> run server

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error")
	}

	client, _, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("db error")
	}

	defer func() {
		if err := db.Disconnect(client); err != nil {
			fmt.Printf("mongo disconnect error: %v", err)
		}
	}()

	router := server.NewRouter()

	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed")
	}
}
