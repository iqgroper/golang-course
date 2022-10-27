package session

import (
	"fmt"
	"net/http"
	"redditclone/pkg/user"
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
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return ExampleTokenSecret, nil
	}

	log.Println("session id value:", sessionCookie.Value)

	token, err := jwt.Parse(sessionCookie.Value, hashSecretGetter)
	if err != nil || !token.Valid {
		log.Println("Error parsing jwt in Check:", err.Error())
		return nil, fmt.Errorf("error parsing jwt in Check or token is invalid: %s", err.Error())
	}

	// payload, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	// 	log.Println("Error getting payload from jwt in Check:", err.Error())
	// }

	sm.mu.RLock()
	sess, ok := sm.data[sessionCookie.Value]
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

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: sess.Expires,
	}
	http.SetCookie(w, cookie)
	return sess, nil
}

// func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
// 	sess, err := SessionFromContext(r.Context())
// 	if err != nil {
// 		return err
// 	}
// 	sm.mu.Lock()
// 	delete(sm.data, sess.ID)
// 	sm.mu.Unlock()
// 	cookie := http.Cookie{
// 		Name:    "session_id",
// 		Expires: time.Now().AddDate(0, 0, -1),
// 		Path:    "/",
// 	}
// 	http.SetCookie(w, &cookie)
// 	return nil
// }
