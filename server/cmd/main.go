package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/noahjonesx/breeze/server/internal/config"
	"github.com/noahjonesx/breeze/server/internal/scheduler"
	"github.com/noahjonesx/breeze/server/internal/store"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sched := scheduler.New(cfg, db)

	log.Println("breeze server starting...")

	if err := sched.Run(); err != nil {
		log.Printf("scheduler error: %v", err)
	}

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := sched.Run(); err != nil {
			log.Printf("scheduler error: %v", err)
		}
	}
}
