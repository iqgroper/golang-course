package posts

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PostsMongoRepository struct {
	DB     *mongo.Collection
	Ctx    *context.Context
	Cancel context.CancelFunc
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
	return &PostsMongoRepository{DB: collection, Ctx: &ctx, Cancel: cancel}
}

func (repo *PostsMongoRepository) GetAll() ([]*Post, error) {
	cur, err := repo.DB.Find(*repo.Ctx, bson.D{})
	if err != nil {
		fmt.Println("ERROR GETTING ALL POSTS", err)
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
	err := repo.DB.FindOne(*repo.Ctx, filter).Decode(&result)
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
		Comments:         make([]*Comment, 0, 10),
		CreatedDTTM:      time.Now(),
		Text:             item.Text,
		URL:              item.URL,
		Type:             item.Type,
		UpvotePercentage: 100,
		Views:            0,
		Author:           AuthorStruct{item.Author.Login, item.Author.ID},
	}

	result, err := repo.DB.InsertOne(*repo.Ctx, newPost)
	if err != nil {
		fmt.Println("ADDING POST", err)
		return nil, err
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		ID := id.Hex()
		newPost.ID = ID
		update := bson.M{
			"$set": bson.M{"ID": ID},
		}
		res := repo.DB.FindOneAndUpdate(*repo.Ctx, bson.M{"_id": id}, update)
		if res.Err() == mongo.ErrNoDocuments {
			fmt.Println("record does not exist")
			return nil, ErrNoPost
		} else if res.Err() != nil {
			log.Fatal("FindOneAndUpdate err", res.Err().Error())
		}
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

	_, err := repo.DB.DeleteOne(*repo.Ctx, filter)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func (repo *PostsMongoRepository) GetByUser(user_login string) ([]*Post, error) {

	filter := bson.M{
		"author": bson.M{
			"username": user_login,
			"id": bson.M{
				"$regex": primitive.Regex{
					Pattern: "[0-9a-z]+",
					Options: "i",
				},
			},
		},
	}

	cur, err := repo.DB.Find(*repo.Ctx, filter)
	if err != nil {
		fmt.Println("ERR GETTING ALL POSTS BY USER", err)
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
			fmt.Println("Error decoding in getting all posts by user")
			return nil, fmt.Errorf("error decoding in getting all posts by user")
		}
		PostList = append(PostList, newPost)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return PostList, nil
}

func (repo *PostsMongoRepository) GetAllByCategory(category string) ([]*Post, error) {

	filter := bson.M{"category": category}

	cur, err := repo.DB.Find(*repo.Ctx, filter)
	if err != nil {
		fmt.Println("ERR GetAllByCategory", err)
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

func (repo *PostsMongoRepository) UpVote(post_id string, username string) (*Post, error) {

	result := &Post{}
	objectId, errGettingObject := primitive.ObjectIDFromHex(post_id)
	if errGettingObject != nil {
		log.Println("Error getting object from id:", post_id)
		return nil, errGettingObject
	}

	filter := bson.M{"_id": objectId}
	err := repo.DB.FindOne(*repo.Ctx, filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if err != nil {
		log.Fatal("find one err", err)
	}

	for _, voter := range result.VotesList {
		if voter.User == username && voter.Vote == 1 {
			return nil, ErrNoCanDo
		}
		if voter.User == username && voter.Vote == -1 {
			_, _ = repo.UnVote(post_id, username)
			result, _ = repo.UpVote(post_id, username)
			return result, nil
		}
	}

	result.Score++
	result.VotesList = append(result.VotesList, VoteStruct{username, 1})
	result.UpvotePercentage = percetageCount(result.VotesList)

	update := bson.M{
		"$set": bson.M{
			"Score":            result.Score,
			"UpvotePercentage": result.UpvotePercentage,
			"VotesList":        result.VotesList,
		},
	}
	res := repo.DB.FindOneAndUpdate(*repo.Ctx, filter, update)

	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if res.Err() != nil {
		log.Fatal(res.Err())
	}

	return result, nil
}

func (repo *PostsMongoRepository) DownVote(post_id string, username string) (*Post, error) {

	result := &Post{}
	objectId, errGettingObject := primitive.ObjectIDFromHex(post_id)
	if errGettingObject != nil {
		log.Println("Error getting object from id:", post_id)
		return nil, errGettingObject
	}

	filter := bson.M{"_id": objectId}
	err := repo.DB.FindOne(*repo.Ctx, filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if err != nil {
		log.Fatal(err)
	}

	for _, voter := range result.VotesList {
		if voter.User == username && voter.Vote == -1 {
			return nil, ErrNoCanDo
		}
		if voter.User == username && voter.Vote == 1 {
			_, _ = repo.UnVote(post_id, username)
			result, _ = repo.DownVote(post_id, username)
			return result, nil
		}
	}

	result.Score--
	result.VotesList = append(result.VotesList, VoteStruct{username, -1})
	result.UpvotePercentage = percetageCount(result.VotesList)

	update := bson.M{
		"$set": bson.M{
			"Score":            result.Score,
			"UpvotePercentage": result.UpvotePercentage,
			"VotesList":        result.VotesList,
		},
	}
	res := repo.DB.FindOneAndUpdate(*repo.Ctx, filter, update)

	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if res.Err() != nil {
		log.Fatal(res.Err())
	}

	return result, nil

}

func (repo *PostsMongoRepository) UnVote(post_id string, username string) (*Post, error) {

	result := &Post{}
	objectId, errGettingObject := primitive.ObjectIDFromHex(post_id)
	if errGettingObject != nil {
		log.Println("Error getting object from id:", post_id)
		return nil, errGettingObject
	}

	filter := bson.M{"_id": objectId}
	err := repo.DB.FindOne(*repo.Ctx, filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if err != nil {
		log.Fatal("FindOne err", err.Error())
	}

	vote := 0
	indexDel := -1
	for idx, voter := range result.VotesList {
		if voter.User == username {
			vote = voter.Vote
			indexDel = idx
			break
		}
	}
	if indexDel == -1 {
		return nil, ErrNoPost
	}

	result.Score -= vote

	if indexDel < len(result.VotesList)-1 {
		copy(result.VotesList[indexDel:], result.VotesList[indexDel+1:])
	}
	result.VotesList[len(result.VotesList)-1] = VoteStruct{}
	result.VotesList = result.VotesList[:len(result.VotesList)-1]

	result.UpvotePercentage = percetageCount(result.VotesList)
	// result.IdMongo = objectId

	update := bson.M{
		"$set": bson.M{
			"Score":            result.Score,
			"UpvotePercentage": result.UpvotePercentage,
			"VotesList":        result.VotesList,
		},
	}
	res := repo.DB.FindOneAndUpdate(*repo.Ctx, filter, update)

	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if res.Err() != nil {
		log.Fatal("FindOneAndUpdate err", res.Err().Error())
	}

	return result, nil

}
