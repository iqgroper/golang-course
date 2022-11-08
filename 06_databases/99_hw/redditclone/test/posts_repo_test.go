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
			{"_id", expectedPosts[0].IdMongo},
			{"ID", expectedPosts[0].ID},
			{"Title", expectedPosts[0].Title},
			{"Score", expectedPosts[0].Score},
			{"VotesList", expectedPosts[0].VotesList},
			{"Category", expectedPosts[0].Category},
			{"Comments", expectedPosts[0].Comments},
			{"CreatedDTTM", expectedPosts[0].CreatedDTTM},
			{"Text", expectedPosts[0].Text},
			{"Type", expectedPosts[0].Type},
			{"UpvotePercentage", expectedPosts[0].UpvotePercentage},
			{"Views", expectedPosts[0].Views},
			{"Author", expectedPosts[0].Author},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedPosts[1].IdMongo},
			{"ID", expectedPosts[1].ID},
			{"Title", expectedPosts[1].Title},
			{"Score", expectedPosts[1].Score},
			{"VotesList", expectedPosts[1].VotesList},
			{"Category", expectedPosts[1].Category},
			{"Comments", expectedPosts[1].Comments},
			{"CreatedDTTM", expectedPosts[1].CreatedDTTM},
			{"Text", expectedPosts[1].Text},
			{"Type", expectedPosts[1].Type},
			{"UpvotePercentage", expectedPosts[1].UpvotePercentage},
			{"Views", expectedPosts[1].Views},
			{"Author", expectedPosts[1].Author},
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

		mt.AddMockResponses(bson.D{{"ok", 0}})

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
			{"_id", expectedPosts[1].IdMongo},
			{"ID", expectedPosts[1].ID},
			{"Title", expectedPosts[1].Title},
			{"Score", expectedPosts[1].Score},
			{"VotesList", expectedPosts[1].VotesList},
			{"Category", expectedPosts[1].Category},
			{"Comments", expectedPosts[1].Comments},
			{"CreatedDTTM", expectedPosts[1].CreatedDTTM},
			{"Text", expectedPosts[1].Text},
			{"Type", expectedPosts[1].Type},
			{"UpvotePercentage", expectedPosts[1].UpvotePercentage},
			{"Views", expectedPosts[1].Views},
			{"Author", expectedPosts[1].Author},
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

		mt.AddMockResponses(bson.D{{"ok", 0}})

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
			{"_id", expectedPosts[0].IdMongo},
			{"ID", expectedPosts[0].ID},
			{"Title", expectedPosts[0].Title},
			{"Score", expectedPosts[0].Score},
			{"VotesList", expectedPosts[0].VotesList},
			{"Category", expectedPosts[0].Category},
			{"Comments", expectedPosts[0].Comments},
			{"CreatedDTTM", expectedPosts[0].CreatedDTTM},
			{"Text", expectedPosts[0].Text},
			{"Type", expectedPosts[0].Type},
			{"UpvotePercentage", expectedPosts[0].UpvotePercentage},
			{"Views", expectedPosts[0].Views},
			{"Author", expectedPosts[0].Author},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedPosts[1].IdMongo},
			{"ID", expectedPosts[1].ID},
			{"Title", expectedPosts[1].Title},
			{"Score", expectedPosts[1].Score},
			{"VotesList", expectedPosts[1].VotesList},
			{"Category", expectedPosts[1].Category},
			{"Comments", expectedPosts[1].Comments},
			{"CreatedDTTM", expectedPosts[1].CreatedDTTM},
			{"Text", expectedPosts[1].Text},
			{"Type", expectedPosts[1].Type},
			{"UpvotePercentage", expectedPosts[1].UpvotePercentage},
			{"Views", expectedPosts[1].Views},
			{"Author", expectedPosts[1].Author},
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

		mt.AddMockResponses(bson.D{{"ok", 0}})

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
			{"_id", expectedPost.IdMongo},
			{"ID", expectedPost.ID},
			{"Title", expectedPost.Title},
			{"Score", expectedPost.Score},
			{"VotesList", expectedPost.VotesList},
			{"Category", expectedPost.Category},
			{"Comments", expectedPost.Comments},
			{"CreatedDTTM", expectedPost.CreatedDTTM},
			{"Text", expectedPost.Text},
			{"Type", expectedPost.Type},
			{"UpvotePercentage", expectedPost.UpvotePercentage},
			{"Views", expectedPost.Views},
			{"Author", expectedPost.Author},
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

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

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

		mt.AddMockResponses(bson.D{{"ok", 0}, {"acknowledged", true}, {"n", 0}})

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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedPost.Score},
					{"UpvotePercentage", expectedPost.UpvotePercentage},
					{"VotesList", expectedPost.VotesList},
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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
			primitive.D{{"ok", 1}},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}},
			),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedAfterUnvotePost.Score},
					{"UpvotePercentage", expectedAfterUnvotePost.UpvotePercentage},
					{"VotesList", expectedAfterUnvotePost.VotesList},
				}},
			},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", expectedAfterUnvotePost.IdMongo},
				{"ID", expectedAfterUnvotePost.ID},
				{"Title", expectedAfterUnvotePost.Title},
				{"Score", expectedAfterUnvotePost.Score},
				{"VotesList", expectedAfterUnvotePost.VotesList},
				{"Category", expectedAfterUnvotePost.Category},
				{"Comments", expectedAfterUnvotePost.Comments},
				{"CreatedDTTM", expectedAfterUnvotePost.CreatedDTTM},
				{"Text", expectedAfterUnvotePost.Text},
				{"Type", expectedAfterUnvotePost.Type},
				{"UpvotePercentage", expectedAfterUnvotePost.UpvotePercentage},
				{"Views", expectedAfterUnvotePost.Views},
				{"Author", expectedAfterUnvotePost.Author}}),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedPost.Score},
					{"UpvotePercentage", expectedPost.UpvotePercentage},
					{"VotesList", expectedPost.VotesList},
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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedPost.Score},
					{"UpvotePercentage", expectedPost.UpvotePercentage},
					{"VotesList", expectedPost.VotesList},
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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
			primitive.D{{"ok", 1}},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}},
			),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedAfterUnvotePost.Score},
					{"UpvotePercentage", expectedAfterUnvotePost.UpvotePercentage},
					{"VotesList", expectedAfterUnvotePost.VotesList},
				}},
			},
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{"_id", expectedAfterUnvotePost.IdMongo},
				{"ID", expectedAfterUnvotePost.ID},
				{"Title", expectedAfterUnvotePost.Title},
				{"Score", expectedAfterUnvotePost.Score},
				{"VotesList", expectedAfterUnvotePost.VotesList},
				{"Category", expectedAfterUnvotePost.Category},
				{"Comments", expectedAfterUnvotePost.Comments},
				{"CreatedDTTM", expectedAfterUnvotePost.CreatedDTTM},
				{"Text", expectedAfterUnvotePost.Text},
				{"Type", expectedAfterUnvotePost.Type},
				{"UpvotePercentage", expectedAfterUnvotePost.UpvotePercentage},
				{"Views", expectedAfterUnvotePost.Views},
				{"Author", expectedAfterUnvotePost.Author}}),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedPost.Score},
					{"UpvotePercentage", expectedPost.UpvotePercentage},
					{"VotesList", expectedPost.VotesList},
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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
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
				{"_id", post.IdMongo},
				{"ID", post.ID},
				{"Title", post.Title},
				{"Score", post.Score},
				{"VotesList", post.VotesList},
				{"Category", post.Category},
				{"Comments", post.Comments},
				{"CreatedDTTM", post.CreatedDTTM},
				{"Text", post.Text},
				{"Type", post.Type},
				{"UpvotePercentage", post.UpvotePercentage},
				{"Views", post.Views},
				{"Author", post.Author}}),
			primitive.D{{"ok", 1}},
			primitive.D{
				{"ok", 1},
				{"value", bson.D{
					{"Score", expectedPost.Score},
					{"UpvotePercentage", expectedPost.UpvotePercentage},
					{"VotesList", expectedPost.VotesList},
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
