package user

import (
	"errors"
	"sync"
)

var (
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("User already exists")
	ErrBadPass    = errors.New("invald password")
)

type UserMemoryRepository struct {
	data   map[string]*User
	LastID uint
	mu     *sync.RWMutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]*User{
			"admin": {
				ID:       0,
				Login:    "admin",
				password: "admin",
			},
		},
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

	repo.LastID++
	newUser := &User{
		ID:       repo.LastID,
		Login:    login,
		password: pass,
	}
	repo.mu.RLock()
	repo.data[login] = newUser
	repo.mu.RUnlock()

	return newUser, nil
}
