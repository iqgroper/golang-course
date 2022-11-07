package tests

import (
	"fmt"
	"redditclone/pkg/user"
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestAuthorize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	var login string = "admin"
	var pass string = "asdfasdf"
	rows := sqlmock.NewRows([]string{"id", "login", "password"})
	expect := []*user.User{
		{
			ID:       "0",
			Login:    "admin",
			Password: "asdfasdf",
		},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Login, item.Password)
	}

	mock.
		ExpectQuery("SELECT id, login, password FROM items WHERE").
		WithArgs(login).
		WillReturnRows(rows)

	repo := &user.UserMysqlRepository{
		DB: db,
	}
	item, err := repo.Authorize(login, pass)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// wrong password error
	row := sqlmock.NewRows([]string{"id", "login", "password"})
	row = row.AddRow("0", "admin", "asdfasdf")

	mock.
		ExpectQuery("SELECT id, login, password FROM items WHERE").
		WithArgs(login).
		WillReturnRows(row)

	_, err = repo.Authorize(login, "wrongpassword")
	fmt.Println("getting", err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// noUser error
	login = "admin111111"
	pass = "asdfasdf"
	blankRows := sqlmock.NewRows([]string{"id", "login", "password"})
	mock.
		ExpectQuery("SELECT id, login, password FROM items WHERE").
		WithArgs(login).
		WillReturnRows(blankRows)

	_, err = repo.Authorize(login, pass)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	//db error
	mock.
		ExpectQuery("SELECT id, login, password FROM items WHERE").
		WithArgs(login).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Authorize(login, pass)
	if errMock := mock.ExpectationsWereMet(); errMock != nil {
		t.Errorf("there were unfulfilled expectations: %s\nReturned: %s", errMock, err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRegister(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var login string = "admin"
	var pass string = "asdfasdf"

	// good query
	rows := sqlmock.NewRows([]string{"id", "login", "password"})
	expect := []*user.User{{
		ID:       "0",
		Login:    "admin",
		Password: "asdfasdf"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Login, item.Password)
	}

	mock.
		ExpectExec("INSERT INTO items").
		WithArgs(login, pass).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &user.UserMysqlRepository{
		DB: db,
	}
	item, err := repo.Register(login, pass)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	//db error
	mock.
		ExpectExec("INSERT INTO items").
		WithArgs(login, pass).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Register(login, pass)
	if errMock := mock.ExpectationsWereMet(); errMock != nil {
		t.Errorf("there were unfulfilled expectations: %s\nReturned: %s", errMock, err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
