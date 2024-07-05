package main

import (
	"log"
	_ "main/docs"
	"main/internal/database"
	task_handlers "main/internal/handlers/task_handlers"
	user_handlers "main/internal/handlers/user_handlers"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Toda App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8080
// @BasePath /

func main() {
	database.InitDatabase()

	// Регистрация обработчиков HTTP маршрутов
	http.HandleFunc("/users", user_handlers.GetUsers)
	http.HandleFunc("/usercreate", user_handlers.CreateUser)
	http.HandleFunc("/userdelete/{id}", user_handlers.DeleteUser)
	http.HandleFunc("/userupdate/{id}", user_handlers.UpdateUser)
	http.HandleFunc("/tasks/start/{id}", task_handlers.StartTask)
	http.HandleFunc("/tasks/end/{id}", task_handlers.EndTask)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Запуск веб-сервера на порту 8080
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
