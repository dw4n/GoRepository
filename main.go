package main

import (
	"gorepository/model"
	"gorepository/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:           "identity front end api",
		EnablePrintRoutes: true,
	})

	dsn := "host=localhost user=postgres password=admin dbname=goblog port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.Post{})

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Update the call to SetupRoutes to pass both repositories
	SetupRoutes(app, userRepo, postRepo)

	app.Listen(":3000")
}
