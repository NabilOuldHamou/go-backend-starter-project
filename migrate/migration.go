package main

import (
	"go-backend-starter-project/initializers"
	"go-backend-starter-project/models"
	"log"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	log.Println("Migrating models to DB...")

	err := initializers.DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Fatalf("Automatic migration has failed: %v", err)
	}

	log.Println("Migration sucessful!")
}
