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

func TestGetAllPosts(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		expectedPosts := []*posts.Post{
			{
				ID:               "636a4200d60d8731dede9fbc",
				Title:            "title",
				Score:            1,
				VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
				Category:         "category",
				Comments:         make([]*posts.Comment, 0, 10),
				CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
				Text:             "text",
				Type:             "text",
				UpvotePercentage: 100,
				Views:            0,
				Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
			},
			{
				ID:               "636a4200d60d8731dede9fbd",
				Title:            "title1",
				Score:            1,
				VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
				Category:         "category1",
				Comments:         make([]*posts.Comment, 0, 10),
				CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
				Text:             "text1",
				Type:             "text",
				UpvotePercentage: 100,
				Views:            0,
				Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
			},
		}

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedPosts[0].IdMongo},
			{Key: "ID", Value: expectedPosts[0].ID},
			{Key: "Title", Value: expectedPosts[0].Title},
			{Key: "Score", Value: expectedPosts[0].Score},
			{Key: "VotesList", Value: expectedPosts[0].VotesList},
			{Key: "Category", Value: expectedPosts[0].Category},
			{Key: "Comments", Value: expectedPosts[0].Comments},
			{Key: "CreatedDTTM", Value: expectedPosts[0].CreatedDTTM},
			{Key: "Text", Value: expectedPosts[0].Text},
			{Key: "Type", Value: expectedPosts[0].Type},
			{Key: "UpvotePercentage", Value: expectedPosts[0].UpvotePercentage},
			{Key: "Views", Value: expectedPosts[0].Views},
			{Key: "Author", Value: expectedPosts[0].Author},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expectedPosts[1].IdMongo},
			{Key: "ID", Value: expectedPosts[1].ID},
			{Key: "Title", Value: expectedPosts[1].Title},
			{Key: "Score", Value: expectedPosts[1].Score},
			{Key: "VotesList", Value: expectedPosts[1].VotesList},
			{Key: "Category", Value: expectedPosts[1].Category},
			{Key: "Comments", Value: expectedPosts[1].Comments},
			{Key: "CreatedDTTM", Value: expectedPosts[1].CreatedDTTM},
			{Key: "Text", Value: expectedPosts[1].Text},
			{Key: "Type", Value: expectedPosts[1].Type},
			{Key: "UpvotePercentage", Value: expectedPosts[1].UpvotePercentage},
			{Key: "Views", Value: expectedPosts[1].Views},
			{Key: "Author", Value: expectedPosts[1].Author},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		foundPosts, err := postsRepo.GetAll()
		require.Nil(t, err)
		require.Equal(t, foundPosts, expectedPosts)
	})

	mt.Run("error returning", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		foundPosts, err := postsRepo.GetAll()
		require.NotNil(t, err)
		require.Nil(t, foundPosts)
	})
}

func TestAddPost(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		newPost := &posts.NewPost{
			Type:     "text",
			Title:    "title1",
			Text:     "text1",
			Category: "category1",
			Author: user.User{
				ID:       "author.id",
				Login:    "author.login",
				Password: "asdfasdf",
			},
		}
		expectedPost := &posts.Post{
			ID:               "636a96165a1c148d1fcf85f2",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "ID", Value: "636a4200d60d8731dede9fbc"},
				}},
			},
		)

		recievedPost, err := postsRepo.Add(newPost)

		require.Nil(t, err)
		require.Equal(t, expectedPost.Author, recievedPost.Author)
		require.Equal(t, expectedPost.Category, recievedPost.Category)
		require.Equal(t, expectedPost.Comments, recievedPost.Comments)
		require.Equal(t, expectedPost.Score, recievedPost.Score)
		require.Equal(t, expectedPost.Text, recievedPost.Text)
		require.Equal(t, expectedPost.Title, recievedPost.Title)
		require.Equal(t, expectedPost.Type, recievedPost.Type)
		require.Equal(t, expectedPost.URL, recievedPost.URL)
		require.Equal(t, expectedPost.UpvotePercentage, recievedPost.UpvotePercentage)
		require.Equal(t, expectedPost.Author.ID, recievedPost.Author.ID)
		require.Equal(t, expectedPost.Author.Username, recievedPost.Author.Username)
	})

	mt.Run("err", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}
		newPost := &posts.NewPost{
			Type:     "text",
			Title:    "title1",
			Text:     "text1",
			Category: "category1",
			Author: user.User{
				ID:       "author.id",
				Login:    "author.login",
				Password: "asdfasdf",
			},
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		recievedPost, err := postsRepo.Add(newPost)

		require.NotNil(t, err)
		require.Nil(t, recievedPost)
	})
}

func TestGetAllPostsByCategory(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		expectedPosts := []*posts.Post{
			{
				ID:               "636a4200d60d8731dede9fbc",
				Title:            "title",
				Score:            1,
				VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
				Category:         "funny",
				Comments:         make([]*posts.Comment, 0, 10),
				CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
				Text:             "text",
				Type:             "text",
				UpvotePercentage: 100,
				Views:            0,
				Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
			},
			{
				ID:               "636a4200d60d8731dede9fbd",
				Title:            "title1",
				Score:            1,
				VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
				Category:         "funny",
				Comments:         make([]*posts.Comment, 0, 10),
				CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
				Text:             "text1",
				Type:             "text",
				UpvotePercentage: 100,
				Views:            0,
				Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
			},
		}

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedPosts[0].IdMongo},
			{Key: "ID", Value: expectedPosts[0].ID},
			{Key: "Title", Value: expectedPosts[0].Title},
			{Key: "Score", Value: expectedPosts[0].Score},
			{Key: "VotesList", Value: expectedPosts[0].VotesList},
			{Key: "Category", Value: expectedPosts[0].Category},
			{Key: "Comments", Value: expectedPosts[0].Comments},
			{Key: "CreatedDTTM", Value: expectedPosts[0].CreatedDTTM},
			{Key: "Text", Value: expectedPosts[0].Text},
			{Key: "Type", Value: expectedPosts[0].Type},
			{Key: "UpvotePercentage", Value: expectedPosts[0].UpvotePercentage},
			{Key: "Views", Value: expectedPosts[0].Views},
			{Key: "Author", Value: expectedPosts[0].Author},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expectedPosts[1].IdMongo},
			{Key: "ID", Value: expectedPosts[1].ID},
			{Key: "Title", Value: expectedPosts[1].Title},
			{Key: "Score", Value: expectedPosts[1].Score},
			{Key: "VotesList", Value: expectedPosts[1].VotesList},
			{Key: "Category", Value: expectedPosts[1].Category},
			{Key: "Comments", Value: expectedPosts[1].Comments},
			{Key: "CreatedDTTM", Value: expectedPosts[1].CreatedDTTM},
			{Key: "Text", Value: expectedPosts[1].Text},
			{Key: "Type", Value: expectedPosts[1].Type},
			{Key: "UpvotePercentage", Value: expectedPosts[1].UpvotePercentage},
			{Key: "Views", Value: expectedPosts[1].Views},
			{Key: "Author", Value: expectedPosts[1].Author},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		foundPosts, err := postsRepo.GetAllByCategory("funny")
		require.Nil(t, err)
		require.Equal(t, foundPosts, expectedPosts)
	})

	mt.Run("error returning", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		foundPosts, err := postsRepo.GetAllByCategory("funny")
		require.NotNil(t, err)
		require.Nil(t, foundPosts)
	})
}

func TestGetPostsByUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		expectedPosts := []*posts.Post{
			{
				ID:               "636a4200d60d8731dede9fbc",
				Title:            "title",
				Score:            1,
				VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
				Category:         "music",
				Comments:         make([]*posts.Comment, 0, 10),
				CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
				Text:             "text",
				Type:             "text",
				UpvotePercentage: 100,
				Views:            0,
				Author:           posts.AuthorStruct{Username: "login", ID: "author.id"},
			},
			{
				ID:               "636a4200d60d8731dede9fbd",
				Title:            "title1",
				Score:            1,
				VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
				Category:         "funny",
				Comments:         make([]*posts.Comment, 0, 10),
				CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
				Text:             "text1",
				Type:             "text",
				UpvotePercentage: 100,
				Views:            0,
				Author:           posts.AuthorStruct{Username: "login", ID: "author.id"},
			},
		}

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedPosts[0].IdMongo},
			{Key: "ID", Value: expectedPosts[0].ID},
			{Key: "Title", Value: expectedPosts[0].Title},
			{Key: "Score", Value: expectedPosts[0].Score},
			{Key: "VotesList", Value: expectedPosts[0].VotesList},
			{Key: "Category", Value: expectedPosts[0].Category},
			{Key: "Comments", Value: expectedPosts[0].Comments},
			{Key: "CreatedDTTM", Value: expectedPosts[0].CreatedDTTM},
			{Key: "Text", Value: expectedPosts[0].Text},
			{Key: "Type", Value: expectedPosts[0].Type},
			{Key: "UpvotePercentage", Value: expectedPosts[0].UpvotePercentage},
			{Key: "Views", Value: expectedPosts[0].Views},
			{Key: "Author", Value: expectedPosts[0].Author},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expectedPosts[1].IdMongo},
			{Key: "ID", Value: expectedPosts[1].ID},
			{Key: "Title", Value: expectedPosts[1].Title},
			{Key: "Score", Value: expectedPosts[1].Score},
			{Key: "VotesList", Value: expectedPosts[1].VotesList},
			{Key: "Category", Value: expectedPosts[1].Category},
			{Key: "Comments", Value: expectedPosts[1].Comments},
			{Key: "CreatedDTTM", Value: expectedPosts[1].CreatedDTTM},
			{Key: "Text", Value: expectedPosts[1].Text},
			{Key: "Type", Value: expectedPosts[1].Type},
			{Key: "UpvotePercentage", Value: expectedPosts[1].UpvotePercentage},
			{Key: "Views", Value: expectedPosts[1].Views},
			{Key: "Author", Value: expectedPosts[1].Author},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		foundPosts, err := postsRepo.GetByUser("login")
		require.Nil(t, err)
		require.Equal(t, foundPosts, expectedPosts)
	})

	mt.Run("error returning", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		foundPosts, err := postsRepo.GetByUser("login")
		require.NotNil(t, err)
		require.Nil(t, foundPosts)
	})
}

func TestGetPostByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
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

		foundPost, err := postsRepo.GetByID("636a4200d60d8731dede9fbc")
		require.Nil(t, err)
		require.Equal(t, foundPost, expectedPost)
	})

	mt.Run("invalid post_id", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		foundPost, err := postsRepo.GetByID("12")

		require.NotNil(t, err)
		require.Nil(t, foundPost)
	})
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

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		ok, err := postsRepo.Delete("636a4200d60d8731dede9fbc")

		require.Nil(t, err)
		require.True(t, ok)
	})

	mt.Run("invalid post_id", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		ok, err := postsRepo.Delete("12")

		require.NotNil(t, err)
		require.False(t, ok)
	})

	mt.Run("error returning", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 0}})

		ok, err := postsRepo.Delete("636a4200d60d8731dede9fbc")
		require.NotNil(t, err)
		require.False(t, ok)
	})
}

func TestUpVote(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            2,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedPost.Score},
					{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
					{Key: "VotesList", Value: expectedPost.VotesList},
				}},
			},
		)

		recievedPost, err := postsRepo.UpVote("636a4200d60d8731dede9fbc", "username")

		require.Nil(t, err)
		require.Equal(t, recievedPost, expectedPost)
	})

	mt.Run("invalid post_id", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		post, err := postsRepo.UpVote("12", "username")

		require.NotNil(t, err)
		require.Nil(t, post)
	})

	mt.Run("redirected to UnVote", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            0,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: -1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            2,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedAfterUnvotePost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}},
			),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedAfterUnvotePost.Score},
					{Key: "UpvotePercentage", Value: expectedAfterUnvotePost.UpvotePercentage},
					{Key: "VotesList", Value: expectedAfterUnvotePost.VotesList},
				}},
			},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: expectedAfterUnvotePost.IdMongo},
				{Key: "ID", Value: expectedAfterUnvotePost.ID},
				{Key: "Title", Value: expectedAfterUnvotePost.Title},
				{Key: "Score", Value: expectedAfterUnvotePost.Score},
				{Key: "VotesList", Value: expectedAfterUnvotePost.VotesList},
				{Key: "Category", Value: expectedAfterUnvotePost.Category},
				{Key: "Comments", Value: expectedAfterUnvotePost.Comments},
				{Key: "CreatedDTTM", Value: expectedAfterUnvotePost.CreatedDTTM},
				{Key: "Text", Value: expectedAfterUnvotePost.Text},
				{Key: "Type", Value: expectedAfterUnvotePost.Type},
				{Key: "UpvotePercentage", Value: expectedAfterUnvotePost.UpvotePercentage},
				{Key: "Views", Value: expectedAfterUnvotePost.Views},
				{Key: "Author", Value: expectedAfterUnvotePost.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedPost.Score},
					{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
					{Key: "VotesList", Value: expectedPost.VotesList},
				}},
			},
		)

		recievedPost, err := postsRepo.UpVote("636a4200d60d8731dede9fbc", "username")

		require.Nil(t, err)
		require.Equal(t, recievedPost, expectedPost)
	})

	mt.Run("noCanDo err returned", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "username", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
			// primitive.D{{"ok", 1}},
		)

		recievedPost, err := postsRepo.UpVote("636a4200d60d8731dede9fbc", "username")

		require.NotNil(t, err)
		require.Nil(t, recievedPost)
	})
}

func TestDownVote(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            0,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: -1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 50,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedPost.Score},
					{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
					{Key: "VotesList", Value: expectedPost.VotesList},
				}},
			},
		)

		recievedPost, err := postsRepo.DownVote("636a4200d60d8731dede9fbc", "username")

		require.Nil(t, err)
		require.Equal(t, expectedPost, recievedPost)
	})

	mt.Run("invalid post_id", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		post, err := postsRepo.DownVote("12", "username")

		require.NotNil(t, err)
		require.Nil(t, post)
	})

	mt.Run("redirected to UnVote", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            2,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            0,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: -1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 50,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedAfterUnvotePost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}},
			),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedAfterUnvotePost.Score},
					{Key: "UpvotePercentage", Value: expectedAfterUnvotePost.UpvotePercentage},
					{Key: "VotesList", Value: expectedAfterUnvotePost.VotesList},
				}},
			},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: expectedAfterUnvotePost.IdMongo},
				{Key: "ID", Value: expectedAfterUnvotePost.ID},
				{Key: "Title", Value: expectedAfterUnvotePost.Title},
				{Key: "Score", Value: expectedAfterUnvotePost.Score},
				{Key: "VotesList", Value: expectedAfterUnvotePost.VotesList},
				{Key: "Category", Value: expectedAfterUnvotePost.Category},
				{Key: "Comments", Value: expectedAfterUnvotePost.Comments},
				{Key: "CreatedDTTM", Value: expectedAfterUnvotePost.CreatedDTTM},
				{Key: "Text", Value: expectedAfterUnvotePost.Text},
				{Key: "Type", Value: expectedAfterUnvotePost.Type},
				{Key: "UpvotePercentage", Value: expectedAfterUnvotePost.UpvotePercentage},
				{Key: "Views", Value: expectedAfterUnvotePost.Views},
				{Key: "Author", Value: expectedAfterUnvotePost.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedPost.Score},
					{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
					{Key: "VotesList", Value: expectedPost.VotesList},
				}},
			},
		)

		recievedPost, err := postsRepo.DownVote("636a4200d60d8731dede9fbc", "username")

		require.Nil(t, err)
		require.Equal(t, recievedPost, expectedPost)
	})

	mt.Run("noCanDo err returned", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            1,
			VotesList:        []posts.VoteStruct{{User: "username", Vote: -1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
		)

		recievedPost, err := postsRepo.DownVote("636a4200d60d8731dede9fbc", "username")

		require.NotNil(t, err)
		require.Nil(t, recievedPost)
	})
}

func TestUnVote(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx, cancel := context.WithCancel(context.Background())

	mt.Run("success", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		id, _ := primitive.ObjectIDFromHex("636a4200d60d8731dede9fbc")
		post := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            0,
			VotesList:        []posts.VoteStruct{{User: "author.login", Vote: 1}, {User: "username", Vote: -1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 50,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}
		expectedPost := &posts.Post{
			IdMongo:          id,
			ID:               "636a4200d60d8731dede9fbc",
			Title:            "title1",
			Score:            -1,
			VotesList:        []posts.VoteStruct{{User: "username", Vote: -1}},
			Category:         "category1",
			Comments:         make([]*posts.Comment, 0, 10),
			CreatedDTTM:      time.Date(2022, time.November, 8, 18, 20, 0, 0, time.UTC),
			Text:             "text1",
			Type:             "text",
			UpvotePercentage: 0,
			Views:            0,
			Author:           posts.AuthorStruct{Username: "author.login", ID: "author.id"},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: post.IdMongo},
				{Key: "ID", Value: post.ID},
				{Key: "Title", Value: post.Title},
				{Key: "Score", Value: post.Score},
				{Key: "VotesList", Value: post.VotesList},
				{Key: "Category", Value: post.Category},
				{Key: "Comments", Value: post.Comments},
				{Key: "CreatedDTTM", Value: post.CreatedDTTM},
				{Key: "Text", Value: post.Text},
				{Key: "Type", Value: post.Type},
				{Key: "UpvotePercentage", Value: post.UpvotePercentage},
				{Key: "Views", Value: post.Views},
				{Key: "Author", Value: post.Author}}),
			primitive.D{{Key: "ok", Value: 1}},
			primitive.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "Score", Value: expectedPost.Score},
					{Key: "UpvotePercentage", Value: expectedPost.UpvotePercentage},
					{Key: "VotesList", Value: expectedPost.VotesList},
				}},
			},
		)

		recievedPost, err := postsRepo.UnVote("636a4200d60d8731dede9fbc", "author.login")

		require.Nil(t, err)
		require.Equal(t, expectedPost, recievedPost)
	})

	mt.Run("no post id err returned", func(mt *mtest.T) {

		postsRepo := &posts.PostsMongoRepository{
			DB:     mt.Coll,
			Ctx:    &ctx,
			Cancel: cancel,
		}

		recievedPost, err := postsRepo.UnVote("23", "username")

		require.NotNil(t, err)
		require.Nil(t, recievedPost)
	})
}
