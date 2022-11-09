package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"redditclone/pkg/handlers"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
)

func newTestRedis() *redismock.ClientMock {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return redismock.NewNiceMock(client)
}

func TestLoginHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockUserRepo(ctrl)
	sessMock := newTestRedis()
	sm := &session.SessionsRedisManager{
		RedisConn: sessMock,
	}

	service := &handlers.UserHandler{
		UserRepo: mockUserRepo,
		Sessions: sm,
		Logger:   log.WithFields(log.Fields{}),
	}

	newUser := &user.User{
		ID:       "0",
		Login:    "asdf",
		Password: "asdfasdf",
	}

	firstUser := &user.NewUser{
		Username: "asdf",
		Password: "asdfasdf",
	}

	//successful
	sess := session.NewSession(newUser)
	dataSerialized, _ := json.Marshal(sess)
	mkey := "sessions:" + sess.ID

	mockUserRepo.EXPECT().Authorize("asdf", "asdfasdf").Return(newUser, nil)
	sessMock.On("Set", mkey, string(dataSerialized), time.Hour).Return(redis.NewStatusResult("OK", nil))

	reqBody, _ := json.Marshal(firstUser)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bytes.NewReader(reqBody))

	service.Login(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	expected := fmt.Sprintf(`{"token": "%s"}`, sess.ID)

	if expected != string(body) {
		t.Errorf("wrong answer, got: %s, expected %ss", expected, string(body))
		return
	}

	resp.Body.Close()

	//auth err no user
	mockUserRepo.EXPECT().Authorize("nosuchuser", "doesnt matter").Return(nil, user.ErrNoUser)
	firstUser = &user.NewUser{
		Username: "nosuchuser",
		Password: "doesnt matter",
	}

	reqBody, _ = json.Marshal(firstUser)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(reqBody))

	service.Login(w, r)

	resp = w.Result()
	if resp.StatusCode != 401 {
		t.Errorf("wrong status code, got: %d, expected 401", resp.StatusCode)
		return
	}

	//auth err bad pass
	mockUserRepo.EXPECT().Authorize("asdf", "wrong").Return(nil, fmt.Errorf("random err"))
	firstUser = &user.NewUser{
		Username: "asdf",
		Password: "wrong",
	}

	reqBody, _ = json.Marshal(firstUser)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(reqBody))

	service.Login(w, r)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("wrong status code, got: %d, expected 500", resp.StatusCode)
		return
	}

	//auth err bad pass
	mockUserRepo.EXPECT().Authorize("asdf", "wrong").Return(nil, user.ErrBadPass)
	firstUser = &user.NewUser{
		Username: "asdf",
		Password: "wrong",
	}

	reqBody, _ = json.Marshal(firstUser)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(reqBody))

	service.Login(w, r)

	resp = w.Result()
	if resp.StatusCode != 401 {
		t.Errorf("wrong status code, got: %d, expected 401", resp.StatusCode)
		return
	}
}

func TestRegisterHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockUserRepo(ctrl)
	sessMock := newTestRedis()
	sm := &session.SessionsRedisManager{
		RedisConn: sessMock,
	}

	service := &handlers.UserHandler{
		UserRepo: mockUserRepo,
		Sessions: sm,
		Logger:   log.WithFields(log.Fields{}),
	}

	newUser := &user.User{
		ID:       "0",
		Login:    "asdf",
		Password: "asdfasdf",
	}

	firstUser := &user.NewUser{
		Username: "asdf",
		Password: "asdfasdf",
	}

	//successful
	sess := session.NewSession(newUser)
	dataSerialized, _ := json.Marshal(sess)
	mkey := "sessions:" + sess.ID

	mockUserRepo.EXPECT().Register("asdf", "asdfasdf").Return(newUser, nil)
	sessMock.On("Set", mkey, string(dataSerialized), time.Hour).Return(redis.NewStatusResult("OK", nil))

	reqBody, _ := json.Marshal(firstUser)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/register", bytes.NewReader(reqBody))

	service.Register(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("wrong status code, got: %d, expected 200", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)

	expected := fmt.Sprintf(`{"token": "%s"}`, sess.ID)

	if expected != string(body) {
		t.Errorf("wrong answer, got: %s, expected %ss", expected, string(body))
		return
	}

	resp.Body.Close()

	//err suer exists
	mockUserRepo.EXPECT().Register("asdf", "asdfasdf").Return(nil, user.ErrUserExists)

	reqBody, _ = json.Marshal(firstUser)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/register", bytes.NewReader(reqBody))

	service.Register(w, r)

	resp = w.Result()
	if resp.StatusCode != 422 {
		t.Errorf("wrong status code, got: %d, expected 422", resp.StatusCode)
		return
	}

}

func TestErrSessionHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockUserRepo(ctrl)
	sessMock := newTestRedis()
	sm := &session.SessionsRedisManager{
		RedisConn: sessMock,
	}

	service := &handlers.UserHandler{
		UserRepo: mockUserRepo,
		Sessions: sm,
		Logger:   log.WithFields(log.Fields{}),
	}

	newUser := &user.User{
		ID:       "0",
		Login:    "asdf",
		Password: "asdfasdf",
	}

	firstUser := &user.NewUser{
		Username: "asdf",
		Password: "asdfasdf",
	}

	sess := session.NewSession(newUser)
	dataSerializedEr, _ := json.Marshal(sess)
	mkeyEr := "sessions:" + sess.ID

	//err getting session login
	mockUserRepo.EXPECT().Authorize("asdf", "asdfasdf").Return(newUser, nil)
	sessMock.On("Set", mkeyEr, string(dataSerializedEr), time.Hour).Return(redis.NewStatusResult("FAIL", fmt.Errorf("some err")))

	reqBodyEr, _ := json.Marshal(firstUser)

	wEr := httptest.NewRecorder()
	rEr := httptest.NewRequest("POST", "/api/login", bytes.NewReader(reqBodyEr))

	service.Login(wEr, rEr)

	respEr := wEr.Result()
	if respEr.StatusCode != 500 {
		t.Errorf("wrong status code, got: %d, expected 500", respEr.StatusCode)
		return
	}

	//err getting session register
	mockUserRepo.EXPECT().Register("asdf", "asdfasdf").Return(newUser, nil)
	sessMock.On("Set", mkeyEr, string(dataSerializedEr), time.Hour).Return(redis.NewStatusResult("FAIL", fmt.Errorf("some err")))

	reqBodyEr, _ = json.Marshal(firstUser)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bytes.NewReader(reqBodyEr))

	service.Register(w, r)

	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("wrong status code, got: %d, expected 500", resp.StatusCode)
		return
	}

}
