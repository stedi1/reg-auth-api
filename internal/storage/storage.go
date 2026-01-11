package storage

import (
	"database/sql"
	_ "embed"
	"log"
)

//go:embed init.sql
var initSQL string

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"-"`
	Email        string `json:"email"`
	RegisterDate string `json:"registerdate"`
}

type Storage struct {
	db *sql.DB
}

// открытие и подготовка хранилища к работе
func NewStorage(db *sql.DB) *Storage {
	// создает таблицу при ее отсутсвии
	initTables(db)
	return &Storage{db: db}
}

// добавление Юзера в таблицу
func (s Storage) AddUser(u User) (id int, err error) {
	res, err := s.db.Exec(`INSERT INTO users 
	(login, password, email, registerdate) VALUES 
	(:login, :password, :email, :registerdate)`,
		sql.Named("login", u.Login),
		sql.Named("password", u.Password),
		sql.Named("email", u.Email),
		sql.Named("registerdate", u.RegisterDate),
	)
	if err != nil {
		log.Println("ошибка добавления пользователя:", err)
		return 0, err
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userId), err
}

// получение юзера по ID
func (s Storage) GetUserByID(userID int) (User, error) {
	u := User{}
	row := s.db.QueryRow(`SELECT id, login, password, email, registerdate FROM
	users WHERE id = :id`, sql.Named("id", userID))
	err := row.Scan(&u.Id, &u.Login, &u.Password, &u.Email, &u.RegisterDate)

	return u, err
}

func initTables(db *sql.DB) {
	_, err := db.Exec(initSQL)
	if err != nil {
		log.Fatal("ошибка создания таблицы", err)
	}
}
