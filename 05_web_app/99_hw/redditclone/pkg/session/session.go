package session

import (
	"context"
	"errors"
	"redditclone/pkg/user"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var (
	ExampleTokenSecret = "secret"
)

type Session struct {
	ID      string
	User    *user.User
	Iat     time.Time
	Expires time.Time
}

func NewSession(user *user.User) *Session {

	iat := time.Now()
	expires := time.Now().Add(90 * 24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expires.Unix(),
		"iat": iat.Unix(),
		"user": map[string]string{
			"username": user.Login,
			"id":       strconv.FormatUint(uint64(user.ID), 10)},
	})

	tokenString, err := token.SignedString([]byte(ExampleTokenSecret))
	if err != nil {
		log.Println("New Session function:", err.Error())
		return nil
	}

	return &Session{
		ID:      tokenString,
		User:    user,
		Iat:     iat,
		Expires: expires,
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
