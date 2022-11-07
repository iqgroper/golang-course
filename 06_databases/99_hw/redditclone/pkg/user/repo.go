package user

import (
	"errors"
	"strconv"
	"sync"
)

var (
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("already exists")
	ErrBadPass    = errors.New("invald password")
)

type UserMemoryRepository struct {
	data   map[string]*User
	LastID int
	mu     *sync.RWMutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]*User{
			"admin": {
				ID:       strconv.Itoa(0),
				Login:    "admin",
				password: "asdfasdf",
			},
		},
		LastID: 0,
		mu:     &sync.RWMutex{},
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	if u.password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *UserMemoryRepository) Register(login, pass string) (*User, error) {
	_, ok := repo.data[login]
	if ok {
		return nil, ErrUserExists
	}

	newUser := &User{
		ID:       strconv.Itoa(repo.LastID),
		Login:    login,
		password: pass,
	}
	repo.LastID++
	repo.mu.RLock()
	repo.data[login] = newUser
	repo.mu.RUnlock()

	return newUser, nil
}
