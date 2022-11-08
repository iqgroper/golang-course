package tests

import (
	"testing"
)

func TestAddComment(t *testing.T) {

	// mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// defer mt.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	// commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)
	// userComment := &user.User{
	// 	ID:       "0",
	// 	Login:    "admin",
	// 	Password: "asdfasdf",
	// }
	// comm := &posts.Comment{
	// 	ID:      "0",
	// 	Body:    "comment",
	// 	Created: time.Now().UTC(),
	// 	Author:  userComment,
	// 	PostID:  "636a4200d60d8731dede9fbc",
	// }

	// mt.Run("success", func(mt *mtest.T) {

	// 	mt.AddMockResponses(mtest.CreateSuccessResponse())

	// 	_, err := commentsRepo.Add(comm.PostID, comm.Body, comm.Author)
	// 	fmt.Println(err)

	// 	assert.Nil(t, err)
	// 	assert.Equal(t, comm, addedComment)
	// })
}

func TestGetAllComments(t *testing.T) {

}

func TestDeleteComment(t *testing.T) {

}
