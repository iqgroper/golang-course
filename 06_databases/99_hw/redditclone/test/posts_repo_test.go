package tests

import (
	"context"
	"redditclone/pkg/posts"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetAllPosts(t *testing.T) {
	// mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// defer mt.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	// postsRepo := &posts.PostsMongoRepository{
	// 	DB:     mt.Coll,
	// 	Ctx:    &ctx,
	// 	Cancel: cancel,
	// }

	// testUser := &user.User{
	// 	ID:       "0",
	// 	Login:    "admin",
	// 	Password: "asdfasdf",
	// }

	// resultPosts := []*posts.Post{
	// 	{
	// 		ID:               "0",
	// 		Title:            "title",
	// 		Score:            1,
	// 		VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
	// 		Category:         "category",
	// 		Comments:         make([]*posts.Comment, 0, 10),
	// 		CreatedDTTM:      time.Now().UTC(),
	// 		Text:             "text",
	// 		Type:             "text",
	// 		UpvotePercentage: 100,
	// 		Views:            0,
	// 		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	// 	},
	// 	{
	// 		ID:               "1",
	// 		Title:            "title1",
	// 		Score:            1,
	// 		VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
	// 		Category:         "category1",
	// 		Comments:         make([]*posts.Comment, 0, 10),
	// 		CreatedDTTM:      time.Now().UTC(),
	// 		Text:             "text1",
	// 		Type:             "text",
	// 		UpvotePercentage: 100,
	// 		Views:            0,
	// 		Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
	// 	},
	// }

	// mt.Run("success", func(mt *mtest.T) {

	// 	mt.AddMockResponses(mtest.CreateSuccessResponse())

	// 	_, err := postsRepo.GetAll()
	// 	fmt.Println(err)

	// 	assert.Nil(t, err)
	// 	assert.Equal(t, comm, addedComment)
	// })
}

func TestAddPost(t *testing.T) {

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

func TestGetAllPostsByCategory(t *testing.T) {

}

func TestGetPostsByUser(t *testing.T) {

}

func TestGetPostByID(t *testing.T) {
	// mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// defer mt.Close()

	// ctx, cancel := context.WithCancel(context.Background())

	// mt.Run("success", func(mt *mtest.T) {

	// 	postsRepo := &posts.PostsMongoRepository{
	// 		DB:     mt.Coll,
	// 		Ctx:    &ctx,
	// 		Cancel: cancel,
	// 	}

	// 	mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
	// 		{"_id", expectedUser.ID},
	// 		{"name", expectedUser.Name},
	// 		{"email", expectedUser.Email},
	// 	}))
	// 	userResponse, err := getFromID(expectedUser.ID)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, &expectedUser, userResponse)
	// })
}

func TestDeletePost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

		ok, err := postsRepo.Delete("636a4200d60d8731dede9fbc")

		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

func TestUpVote(t *testing.T) {

}

func TestDownVote(t *testing.T) {

}

func TestUnVote(t *testing.T) {

}
