package session

import (
	"context"
	"errors"
	"redditclone/pkg/user"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var (
	ExampleTokenSecret = "secret"
)

type Session struct {
	ID   string
	User *user.User
}

func NewSession(user *user.User) *Session {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{
			"username": user.Login,
			"id":       user.ID,
		},
	})

	tokenString, err := token.SignedString([]byte(ExampleTokenSecret))
	if err != nil {
		log.Println("New Session function:", err.Error())
		return nil
	}

	return &Session{
		ID:   tokenString,
		User: user,
	}
}

var (
	ErrNoAuth = errors.New("no session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
