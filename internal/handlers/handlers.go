package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reg-auth-api/internal/services"
	"reg-auth-api/internal/storage"
	"time"
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

}

func (h Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	regReq := RegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := storage.User{
		Login:        regReq.Login,
		Email:        regReq.Email,
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
