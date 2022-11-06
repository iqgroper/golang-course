package posts

import (
	"context"
	"fmt"
	"log"
	"redditclone/pkg/comments"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PostsMongoRepository struct {
	DB     *mongo.Collection
	ctx    *context.Context
	cancel context.CancelFunc
}

func NewMongoRepository() *PostsMongoRepository {
	client, errMongo := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost"))
	if errMongo != nil {
		panic(errMongo)
	}
	ctx, cancel := context.WithCancel(context.Background())
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	errMongo = client.Ping(ctx, readpref.Primary())
	if errMongo != nil {
		panic(errMongo)

	}

	collection := client.Database("coursera").Collection("posts")
	return &PostsMongoRepository{DB: collection, ctx: &ctx, cancel: cancel}
}

func (repo *PostsMongoRepository) GetAll() ([]*Post, error) {
	cur, err := repo.DB.Find(*repo.ctx, bson.D{})
	if err != nil {
		fmt.Println("GETTING ALL POSTS", err)
		return nil, err
	}

	PostList := make([]*Post, 0, 5)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		newPost := &Post{}
		err := cur.Decode(newPost)
		if err != nil {
			fmt.Println("Error decoding in GetAllPosts")
			return nil, fmt.Errorf("error decoding in GetAllPosts")
		}
		PostList = append(PostList, newPost)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return PostList, nil
}

func (repo *PostsMongoRepository) GetByID(id string) (*Post, error) {
	result := &Post{}
	objectId, errGettingObject := primitive.ObjectIDFromHex(id)
	if errGettingObject != nil {
		log.Println("Error getting object from id:", id)
		return nil, errGettingObject
	}
	filter := bson.M{"_id": objectId}
	err := repo.DB.FindOne(*repo.ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func (repo *PostsMongoRepository) Add(item *NewPost) (*Post, error) {

	newPost := &Post{
		Title:            item.Title,
		Score:            1,
		VotesList:        []VoteStruct{{item.Author.Login, 1}},
		Category:         item.Category,
		Comments:         make([]*comments.Comment, 0, 10),
		CreatedDTTM:      time.Now(),
		Text:             item.Text,
		URL:              item.URL,
		Type:             item.Type,
		UpvotePercentage: 100,
		Views:            0,
		Author:           AuthorStruct{item.Author.Login, item.Author.ID},
	}

	result, err := repo.DB.InsertOne(*repo.ctx, newPost)
	if err != nil {
		fmt.Println("ADDING POST", err)
		return nil, err
	}
	// return newPost, nil

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		newPost.ID = id.Hex()
		return newPost, nil
	} else {
		return nil, fmt.Errorf("id assertion failed")
	}
}

func (repo *PostsMongoRepository) Delete(id string) (bool, error) {
	objectId, errGettingObject := primitive.ObjectIDFromHex(id)
	if errGettingObject != nil {
		log.Println("Error getting object from id")
		return false, errGettingObject
	}
	filter := bson.M{"_id": objectId}

	_, err := repo.DB.DeleteOne(*repo.ctx, filter)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func (repo *PostsMongoRepository) GetByUser(user_login string) ([]*Post, error) {
	cur, err := repo.DB.Find(*repo.ctx, bson.M{"Author": bson.M{"User": user_login}})
	if err != nil {
		fmt.Println("GETTING ALL POSTS BY USER", err)
		return nil, err
	}

	PostList := make([]*Post, 0, 5)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		newPost := &Post{}
		err := cur.Decode(newPost)
		if err != nil {
			fmt.Println("Error decoding in GetAllPosts")
			return nil, fmt.Errorf("error decoding in GetAllPosts")
		}
		PostList = append(PostList, newPost)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return PostList, nil
}

func (repo *PostsMongoRepository) GetAllByCategory(category string) ([]*Post, error) {
	// result := make([]*Post, 0, 10)
	// for _, item := range repo.data {
	// 	if item.Category == category {
	// 		result = append(result, item)
	// 	}
	// }

	// if len(result) == 0 {
	// 	return nil, ErrNoPost
	// }

	// return result, nil
	return nil, nil

}

func (repo *PostsMongoRepository) UpVote(post_id string, username string) (*Post, error) {
	// for indexPost, item := range repo.data {
	// 	if item.ID == post_id {
	// 		for _, voter := range item.VotesList {
	// 			if voter.User == username && voter.Vote == 1 {
	// 				return nil, ErrNoCanDo
	// 			}
	// 			if voter.User == username && voter.Vote == -1 {
	// 				repo.data[indexPost], _ = repo.UnVote(post_id, username)
	// 				repo.data[indexPost], _ = repo.UpVote(post_id, username)
	// 				return repo.data[indexPost], nil
	// 			}
	// 		}

	// 		item.Score += 1

	// 		item.VotesList = append(item.VotesList, VoteStruct{username, 1})
	// 		item.UpvotePercentage = percetageCount(item.VotesList)

	// 		return item, nil
	// 	}
	// }
	// return nil, ErrNoPost
	return nil, nil

}

func (repo *PostsMongoRepository) DownVote(post_id string, username string) (*Post, error) {
	// for indexPost, item := range repo.data {
	// 	if item.ID == post_id {
	// 		for _, voter := range item.VotesList {
	// 			if voter.User == username && voter.Vote == -1 {
	// 				return nil, ErrNoCanDo
	// 			}
	// 			if voter.User == username && voter.Vote == 1 {
	// 				repo.data[indexPost], _ = repo.UnVote(post_id, username)
	// 				repo.data[indexPost], _ = repo.DownVote(post_id, username)
	// 				return repo.data[indexPost], nil
	// 			}
	// 		}
	// 		item.Score -= 1

	// 		item.VotesList = append(item.VotesList, VoteStruct{username, -1})
	// 		item.UpvotePercentage = percetageCount(item.VotesList)

	// 		return item, nil
	// 	}
	// }
	// return nil, ErrNoPost
	return nil, nil

}

func (repo *PostsMongoRepository) UnVote(post_id string, username string) (*Post, error) {
	// 	postIndexToRemove := -1
	// 	voteIndexToRemove := -1
	// LOOP:
	// 	for postIdx, item := range repo.data {
	// 		if item.ID == post_id {
	// 			postIndexToRemove = postIdx
	// 			for idx, voter := range item.VotesList {
	// 				if voter.User == username {
	// 					voteIndexToRemove = idx

	// 					item.Score -= voter.Vote

	// 					break LOOP
	// 				}
	// 			}
	// 		}
	// 	}

	// 	repo.mu.Lock()
	// 	if voteIndexToRemove < len(repo.data[postIndexToRemove].VotesList)-1 {
	// 		copy(repo.data[postIndexToRemove].VotesList[voteIndexToRemove:], repo.data[postIndexToRemove].VotesList[voteIndexToRemove+1:])
	// 	}
	// 	repo.data[postIndexToRemove].VotesList[len(repo.data[postIndexToRemove].VotesList)-1] = VoteStruct{}
	// 	repo.data[postIndexToRemove].VotesList = repo.data[postIndexToRemove].VotesList[:len(repo.data[postIndexToRemove].VotesList)-1]
	// 	repo.mu.Unlock()

	// 	repo.data[postIndexToRemove].UpvotePercentage = percetageCount(repo.data[postIndexToRemove].VotesList)

	// return repo.data[postIndexToRemove], nil
	return nil, nil

}
