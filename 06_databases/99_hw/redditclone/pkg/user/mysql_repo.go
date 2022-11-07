package user

import (
	"database/sql"
	"fmt"
	"strconv"
)

type UserMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo() *UserMysqlRepository {
	// основные настройки к базе
	dsn := "root:love@tcp(localhost:3306)/golang?"
	// указываем кодировку
	dsn += "charset=utf8"
	// отказываемся от prapared statements
	// параметры подставляются сразу
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}
	return &UserMysqlRepository{DB: db}
}

func (repo *UserMysqlRepository) Authorize(login, pass string) (*User, error) {
	user := &User{}
	err := repo.DB.QueryRow("SELECT id, login, password FROM items WHERE login = ?;", login).
		Scan(&user.ID, &user.Login, &user.password)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}
	if err != nil {
		fmt.Println("QueryRow", err)
		return nil, err
	}

	if user.password != pass {
		return nil, ErrBadPass
	}

	return user, nil
}

func (repo *UserMysqlRepository) Register(login, pass string) (*User, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO items (`login`, `password`) VALUES (?, ?)",
		login,
		pass,
	)
	if err != nil {
		fmt.Println("QueryRow", err)
		return nil, err
	}

	id, _ := result.LastInsertId()

	newUser := &User{
		ID:       strconv.Itoa(int(id)),
		Login:    login,
		password: pass,
	}

	return newUser, nil
}
