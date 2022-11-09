package tests

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
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

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// mockUserRepo := user.NewMockUserRepo(ctrl)
	// sessMock := newTestRedis()
	// sm := &session.SessionsRedisManager{
	// 	RedisConn: sessMock,
	// }

	// service := &handlers.UserHandler{
	// 	UserRepo: mockUserRepo,
	// 	Sessions: sm,
	// 	Logger:   log.WithFields(log.Fields{}),
	// }

	// newUser := &user.User{
	// 	ID:       "0",
	// 	Login:    "asdf",
	// 	Password: "asdfasdf",
	// }

	// firstUser := &user.NewUser{
	// 	Username: "asdf",
	// 	Password: "asdfasdf",
	// }

	// //successful
	// sess := session.NewSession(newUser)
	// dataSerialized, _ := json.Marshal(sess)
	// mkey := "sessions:" + sess.ID

	// mockUserRepo.EXPECT().Authorize("asdf", "asdfasdf").Return(newUser, nil)
	// sessMock.On("Set", mkey, dataSerialized, 24*time.Hour).Return(redis.NewCmdResult("OK", nil))

	// reqBody, _ := json.Marshal(firstUser)

	// w := httptest.NewRecorder()
	// r := httptest.NewRequest("POST", "/api/posts", bytes.NewReader(reqBody))

	// service.Login(w, r)

}
