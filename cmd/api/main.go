package main

import (
	"learning-platform/internal/configs"
	"learning-platform/internal/di"
	"learning-platform/internal/platform/db"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	if err := db.Migration(cfg); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	pool, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer pool.Close()

	container, err := di.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	r := container.SetupRouter()
	err = r.Run(":" + cfg.AppConfig.AppPort)
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
