package main

import (
	"database/sql"
	"log"
	"reg-auth-api/internal/handlers"
	"reg-auth-api/internal/server"
	"reg-auth-api/internal/services"
	"reg-auth-api/internal/storage"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	// подключаем БД
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Fatal(err, "не удалось создать/открыть БД")
	}
	defer db.Close()

	// создаем Service, который будет работать с юзерами из БД в хендлерах
	s := storage.NewStorage(db)
	UserService := services.NewUserService(s)

	// объект со всеми хендлерами
	handler := handlers.InitHandler(UserService)

	// создаем роутер и настраиваем его
	r := chi.NewRouter()
	r.Get("/", handler.MainPage)
	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)
	r.Get("/get/user/{id}", handler.GetUserByID)

	// создаем и запускаем сервер с нужным нам роутером
	server := server.NewServer(r)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
