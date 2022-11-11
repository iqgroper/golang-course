package user

import (
	"database/sql"
	"strconv"

	"github.com/pkg/errors"
)

type UserMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo() *UserMysqlRepository {
	dsn := "root:love@tcp(localhost:3306)/golang?"
	dsn += "charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &UserMysqlRepository{DB: db}
}

func (repo *UserMysqlRepository) Authorize(login, pass string) (*User, error) {
	user := &User{}
	err := repo.DB.QueryRow("SELECT id, login, password FROM items WHERE login = ?;", login).
		Scan(&user.ID, &user.Login, &user.Password)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}
	if err != nil {
		return nil, errors.Wrap(err, "QueryRowError")
	}

	if user.Password != pass {
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
		return nil, errors.Wrap(err, "QueryRowError")
	}

	id, _ := result.LastInsertId()

	newUser := &User{
		ID:       strconv.Itoa(int(id)),
		Login:    login,
		Password: pass,
	}

	return newUser, nil
}
