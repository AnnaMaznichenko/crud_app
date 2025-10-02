package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"crud_app/api"
	"crud_app/config"
	"crud_app/repository"
	"crud_app/service"
)

func main() {
	fmt.Println("Server starting on :8080")

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	var userRepo repository.UserRepo
	userRepo = repository.NewUserRepo(db)

	var userValidator service.UserValidator
	userValidator = service.NewUserValidator(userRepo)

	var userService service.User
	userService = service.NewUser(userValidator, userRepo)

	r := chi.NewRouter()

	api.SetUserHandlers(r, userService)

	http.ListenAndServe(":8080", r)
}
