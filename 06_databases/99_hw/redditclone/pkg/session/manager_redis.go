package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/pkg/user"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type SessionsRedisManager struct {
	RedisConn redis.Cmdable
}

func NewSessionManager() *SessionsRedisManager {

	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ping := conn.Ping()
	if ping.Err() != nil {
		panic(ping.Err())
	}

	log.Println("Connected to Redis")

	return &SessionsRedisManager{
		RedisConn: conn,
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
		return nil, fmt.Errorf("error parsing jwt in Check or token is invalid: %w", err)
	}

	mkey := "sessions:" + authToken
	data, err := sm.RedisConn.Get(mkey).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "cant get data from Redis in Check")
	}
	sess := &Session{}
	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, errors.Wrap(err, "cant unpack session data in Check")
	}

	return sess, nil
}

func (sm *SessionsRedisManager) Create(w http.ResponseWriter, user *user.User) (*Session, error) {
	sess := NewSession(user)

	dataSerialized, _ := json.Marshal(sess)
	mkey := "sessions:" + sess.ID
	result := sm.RedisConn.Set(mkey, string(dataSerialized), time.Hour)

	res, err := result.Result()

	if res != "OK" || err != nil {
		fmt.Println("error creating session, res:", res)
		return nil, fmt.Errorf("error creating session")
	}

	return sess, nil
}
