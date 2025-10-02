package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"crud_app/dto"
	"crud_app/service"
)

func SetUserHandlers(router *chi.Mux, userService service.User) {
	userRouter := chi.NewRouter()

	userRouter.Get("/list", listUserHandler(userService))

	userRouter.Post("/create", createUserHandler(userService))

	userRouter.Put("/update/{id}", updateUserHandler(userService))

	userRouter.Delete("/delete/{id}", deleteUserHandler(userService))

	router.Mount("/users", userRouter)
}

func listUserHandler(userService service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var result Result
		result.Data, result.Error = userService.List(ctx)

		writeResponseWithJson(w, result)
	}
}

func createUserHandler(userService service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var user dto.User
		var result Result

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			result.Error = fmt.Errorf("invalid JSON format")
		} else {
			result.Data, result.Error = userService.Create(ctx, &user)
		}

		writeResponseWithJson(w, result)
	}
}

func updateUserHandler(userService service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		uuid, err := strconv.ParseUint(id, 10, 64)

		var user dto.User
		var result Result

		if err != nil {
			result.Error = fmt.Errorf("id is not uuid")
		} else if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			result.Error = fmt.Errorf("invalid JSON format")
		} else {
			result.Error = userService.Update(ctx, &user, uint(uuid))
		}

		writeResponse(w, result)
	}
}

func deleteUserHandler(userService service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		uuid, err := strconv.ParseUint(id, 10, 64)

		var result Result

		if err != nil {
			result.Error = fmt.Errorf("id is not uuid")
		} else {
			result.Error = userService.Delete(ctx, uint(uuid))
		}

		writeResponse(w, result)
	}
}
