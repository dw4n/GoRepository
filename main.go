package main

import (
	"os"

	"gorepository/model"
	"gorepository/repository"
	"gorepository/routes"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables")
	}

	app := fiber.New(fiber.Config{
		AppName:           "identity front end api",
		EnablePrintRoutes: true,
	})

	// Retrieve connection string from environment variables
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.Post{})

	repos := repository.NewRepositories(db)

	// Update the call to SetupRoutes to pass both repositories
	routes.SetupRoutes(app, repos)

	app.Listen(":3000")
}
