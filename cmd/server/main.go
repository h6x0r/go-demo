package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/h6x0r/go-demo/internal/database"
	"github.com/h6x0r/go-demo/internal/handlers"
	"github.com/h6x0r/go-demo/internal/repository"
)

func main() {
	// Подключаемся к базе данных
	database.Connect()
	defer database.DB.Close()

	// Выполняем миграцию (создаём таблицы, если их ещё нет)
	database.Migrate()

	repo := repository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(repo)

	r := gin.Default()

	r.POST("/user", userHandler.CreateUser)
	r.GET("/user/:id", userHandler.GetUser)
	r.PATCH("/user/:id", userHandler.UpdateUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)

	log.Println("Сервер запущен на порту 8080")
	r.Run(":8080")
}
