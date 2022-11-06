package session

import (
	"fmt"
	"net/http"
	"redditclone/pkg/user"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type SessionsManager struct {
	data map[string]*Session
	mu   *sync.RWMutex
}

func NewSessionsManager() *SessionsManager {
	return &SessionsManager{
		data: make(map[string]*Session, 10),
		mu:   &sync.RWMutex{},
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {

	authTokenStr, ok := r.Header["Authorization"]
	if !ok {
		return nil, ErrNoAuth
	}

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return []byte(ExampleTokenSecret), nil
	}

	authToken := strings.Split(authTokenStr[0], " ")[1]
	token, err := jwt.Parse(authToken, hashSecretGetter)
	if err != nil || !token.Valid {
		log.Println("Error parsing jwt in Check:", err.Error())
		return nil, fmt.Errorf("error parsing jwt in Check or token is invalid: %s", err.Error())
	}

	sm.mu.RLock()
	sess, ok := sm.data[authToken]
	sm.mu.RUnlock()

	if !ok {
		return nil, ErrNoAuth
	}

	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, user *user.User) (*Session, error) {
	sess := NewSession(user)

	sm.mu.Lock()
	sm.data[sess.ID] = sess
	sm.mu.Unlock()

	return sess, nil
}