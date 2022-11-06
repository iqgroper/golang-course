package user

import (
	"database/sql"
	"fmt"
)

type UserMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo(db *sql.DB) *UserMysqlRepository {
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
		ID:       uint(id),
		Login:    login,
		password: pass,
	}

	return newUser, nil
}
