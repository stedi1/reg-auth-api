package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reg-auth-api/internal/services"
	"reg-auth-api/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	UserService *services.UserService
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func InitHandler(s *services.UserService) *Handler {
	return &Handler{UserService: s}
}

func (h Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok, сервер запущен"))
}

func (h Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.UserService.Storage.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "пользователь не найден: "+err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	err = enc.Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	regReq := RegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := storage.User{
		Login:        strings.ToLower(regReq.Login),
		Email:        strings.ToLower(regReq.Email),
		RegisterDate: time.Now().Format("2006.01.02"),
	}
	// проверяем пароль и хешируем
	hashed, err := h.UserService.CheckAndHashPass(regReq.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Password = hashed
	// проверяем email
	if !h.UserService.IsEmailValid(user.Email) {
		http.Error(w, "incorrect email", http.StatusBadRequest)
		return
	}
	log.Println(user)

	// сохраняем пользователя в БД
	id, err := h.UserService.Storage.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)

}

func (h Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("в процессе"))
}
