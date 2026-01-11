package services

import (
	"errors"
	"reg-auth-api/internal/storage"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Storage *storage.Storage
}

func NewUserService(s *storage.Storage) *UserService {
	return &UserService{Storage: s}
}

// регистрирует пользователя
func (s *UserService) Register(u storage.User) (userID int, err error) {
	id, err := s.Storage.AddUser(u)

	return id, err
}

// првоеряет длину пароля пользователя и хеширует если ок
func (s *UserService) CheckAndHashPass(pass string) (hashedPass string, err error) {
	if len([]rune(pass)) < 8 {
		return "", errors.New("password less than 8 symbols")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), err
}

// проверяет email на корректность
func (s *UserService) IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
