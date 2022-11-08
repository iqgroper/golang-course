package tests

import (
	"context"
	"redditclone/pkg/posts"
	"redditclone/pkg/user"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAddComment(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {
		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		userComment := &user.User{
			ID:       "0",
			Login:    "admin",
			Password: "asdfasdf",
		}

		comm := &posts.Comment{
			ID:      "0",
			Body:    "comment",
			Created: time.Now().UTC(),
			Author:  userComment,
			PostID:  "636a4200d60d8731dede9fbc",
		}
		expectedComments := []*posts.Comment{
			{
				ID:      "0",
				Body:    "comment",
				Created: time.Now().UTC(),
				Author:  userComment,
				PostID:  "636a4200d60d8731dede9fbc",
			},
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         expectedComments,
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: expectedPost.IdMongo},
				{Key: "ID", Value: expectedPost.ID},
				{Key: "Title", Value: expectedPost.Title},
				{Key: "Score", Value: expectedPost.Score},
				{Key: "VotesList", Value: expectedPost.VotesList},
				{Key: "Category", Value: expectedPost.Category},
				{Key: "Comments", Value: expectedPost.Comments},
				{Key: "CreatedDTTM", Value: expectedPost.CreatedDTTM},
				{Key: "Text", Value: expectedPost.Text},
				{Key: "Type", Value: expectedPost.Type},
				{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
				{Key: "Views", Value: expectedPost.Views},
				{Key: "Author", Value: expectedPost.Author},
			}),
			primitive.D{
				{Key: "ok", Value: 1},
			},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Comments", Value: expectedPost.Comments}},
				},
			},
		)

		addedComment, err := commentsRepo.Add(comm.PostID, comm.Body, comm.Author)

		require.Nil(t, err)
		require.Equal(t, comm.Body, addedComment.Body)
		require.Equal(t, comm.ID, addedComment.ID)
		require.Equal(t, comm.PostID, addedComment.PostID)
		require.Equal(t, comm.Author.ID, addedComment.Author.ID)
	})

	mt.Run("err", func(mt *mtest.T) {

		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		userComment := &user.User{
			ID:       "0",
			Login:    "admin",
			Password: "asdfasdf",
		}
		comm := &posts.Comment{
			ID:      "0",
			Body:    "comment",
			Created: time.Now().UTC(),
			Author:  userComment,
			PostID:  "636a4200d60d8731dede9fbc",
		}

		comment, err := commentsRepo.Add("wrong id", comm.Body, comm.Author)

		require.NotNil(t, err)
		require.Nil(t, comment)

	})
}

func TestGetAllComments(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		userComment := &user.User{
			ID:       "0",
			Login:    "admin",
			Password: "asdfasdf",
		}
		expectedComments := []*posts.Comment{
			{
				ID:      "0",
				Body:    "comment",
				Created: time.Now().UTC(),
				Author:  userComment,
				PostID:  "636a4200d60d8731dede9fbc",
			},
			{
				ID:      "1",
				Body:    "comment1",
				Created: time.Now().UTC(),
				Author:  userComment,
				PostID:  "636a4200d60d8731dede9fbc",
			},
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         expectedComments,
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: expectedPost.IdMongo},
				{Key: "ID", Value: expectedPost.ID},
				{Key: "Title", Value: expectedPost.Title},
				{Key: "Score", Value: expectedPost.Score},
				{Key: "VotesList", Value: expectedPost.VotesList},
				{Key: "Category", Value: expectedPost.Category},
				{Key: "Comments", Value: expectedPost.Comments},
				{Key: "CreatedDTTM", Value: expectedPost.CreatedDTTM},
				{Key: "Text", Value: expectedPost.Text},
				{Key: "Type", Value: expectedPost.Type},
				{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
				{Key: "Views", Value: expectedPost.Views},
				{Key: "Author", Value: expectedPost.Author},
			}))

		comments, err := commentsRepo.GetAll("636a4200d60d8731dede9fbc")

		require.Nil(t, err)
		require.Equal(t, expectedComments[0].Author.ID, comments[0].Author.ID)
		require.Equal(t, expectedComments[0].Author.Login, comments[0].Author.Login)
		require.Equal(t, expectedComments[0].Body, comments[0].Body)
		require.Equal(t, expectedComments[0].ID, comments[0].ID)
		require.Equal(t, expectedComments[0].PostID, comments[0].PostID)
		require.Equal(t, expectedComments[1].Author.ID, comments[1].Author.ID)
		require.Equal(t, expectedComments[1].Author.Login, comments[1].Author.Login)
		require.Equal(t, expectedComments[1].Body, comments[1].Body)
		require.Equal(t, expectedComments[1].ID, comments[1].ID)
		require.Equal(t, expectedComments[1].PostID, comments[1].PostID)
	})

	mt.Run("err", func(mt *mtest.T) {

		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		comments, err := commentsRepo.GetAll("23")

		require.NotNil(t, err)
		require.Nil(t, comments)

	})
}

func TestDeleteComment(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {
		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		userComment := &user.User{
			ID:       "0",
			Login:    "admin",
			Password: "asdfasdf",
		}

		expectedComments := []*posts.Comment{
			{
				ID:      "0",
				Body:    "comment",
				Created: time.Now().UTC(),
				Author:  userComment,
				PostID:  "636a4200d60d8731dede9fbc",
			},
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         expectedComments,
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: expectedPost.IdMongo},
				{Key: "ID", Value: expectedPost.ID},
				{Key: "Title", Value: expectedPost.Title},
				{Key: "Score", Value: expectedPost.Score},
				{Key: "VotesList", Value: expectedPost.VotesList},
				{Key: "Category", Value: expectedPost.Category},
				{Key: "Comments", Value: expectedPost.Comments},
				{Key: "CreatedDTTM", Value: expectedPost.CreatedDTTM},
				{Key: "Text", Value: expectedPost.Text},
				{Key: "Type", Value: expectedPost.Type},
				{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
				{Key: "Views", Value: expectedPost.Views},
				{Key: "Author", Value: expectedPost.Author},
			}),
			primitive.D{
				{Key: "ok", Value: 1},
			},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Comments", Value: expectedPost.Comments}},
				},
			},
		)

		ok, err := commentsRepo.Delete("636a4200d60d8731dede9fbc", "0")

		require.Nil(t, err)
		require.True(t, ok)
	})

	mt.Run("success not last comment", func(mt *mtest.T) {
		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		userComment := &user.User{
			ID:       "0",
			Login:    "admin",
			Password: "asdfasdf",
		}

		expectedComments := []*posts.Comment{
			{
				ID:      "0",
				Body:    "comment",
				Created: time.Now().UTC(),
				Author:  userComment,
				PostID:  "636a4200d60d8731dede9fbc",
			},
			{
				ID:      "1",
				Body:    "comment",
				Created: time.Now().UTC(),
				Author:  userComment,
				PostID:  "636a4200d60d8731dede9fbc",
			},
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         expectedComments,
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: expectedPost.IdMongo},
				{Key: "ID", Value: expectedPost.ID},
				{Key: "Title", Value: expectedPost.Title},
				{Key: "Score", Value: expectedPost.Score},
				{Key: "VotesList", Value: expectedPost.VotesList},
				{Key: "Category", Value: expectedPost.Category},
				{Key: "Comments", Value: expectedPost.Comments},
				{Key: "CreatedDTTM", Value: expectedPost.CreatedDTTM},
				{Key: "Text", Value: expectedPost.Text},
				{Key: "Type", Value: expectedPost.Type},
				{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
				{Key: "Views", Value: expectedPost.Views},
				{Key: "Author", Value: expectedPost.Author},
			}),
			primitive.D{
				{Key: "ok", Value: 1},
			},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Comments", Value: expectedPost.Comments}},
				},
			},
		)

		ok, err := commentsRepo.Delete("636a4200d60d8731dede9fbc", "0")

		require.Nil(t, err)
		require.True(t, ok)
	})

	mt.Run("err", func(mt *mtest.T) {

		commentsRepo := posts.NewMongoRepo(mt.Coll, &ctx, cancel)

		ok, err := commentsRepo.Delete("e", "0")

		require.NotNil(t, err)
		require.False(t, ok)

	})
}
