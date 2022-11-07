package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/pkg/user"
	"strings"

	"github.com/gomodule/redigo/redis"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type SessionsRedisManager struct {
	redisConn redis.Conn
}

func NewSessionManager() *SessionsRedisManager {

	conn, err := redis.DialURL("redis://user:@localhost:6379/0")
	if err != nil {
		log.Fatalf("cant connect to redis")
	}

	log.Println("Connected to Redis")

	return &SessionsRedisManager{
		redisConn: conn,
	}
}

func (sm *SessionsRedisManager) Check(r *http.Request) (*Session, error) {

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

	mkey := "sessions:" + authToken
	data, err := redis.Bytes(sm.redisConn.Do("GET", mkey))
	if err != nil {
		log.Println("cant get data:", err)
		return nil, err
	}
	sess := &Session{}
	err = json.Unmarshal(data, sess)
	if err != nil {
		log.Println("cant unpack session data:", err)
		return nil, err
	}

	return sess, nil
}

func (sm *SessionsRedisManager) Create(w http.ResponseWriter, user *user.User) (*Session, error) {
	sess := NewSession(user)

	dataSerialized, _ := json.Marshal(sess)
	mkey := "sessions:" + sess.ID
	result, err := redis.String(sm.redisConn.Do("SET", mkey, dataSerialized, "EX", 86400))
	if err != nil {
		log.Println("Error creating session in redis")
		return nil, err
	}
	if result != "OK" {
		return nil, fmt.Errorf("error creating session")
	}

	return sess, nil
}
