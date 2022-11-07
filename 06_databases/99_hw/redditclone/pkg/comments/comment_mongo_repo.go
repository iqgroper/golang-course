package comments

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"redditclone/pkg/posts"
// 	"redditclone/pkg/user"
// 	"strconv"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"gopkg.in/mgo.v2/bson"
// )

// type CommentMongoRepository struct {
// 	DB     *mongo.Collection
// 	ctx    *context.Context
// 	Cancel context.CancelFunc
// }

// func NewMongoRepo(db *mongo.Collection, ctx *context.Context, cancel context.CancelFunc) *CommentMongoRepository {
// 	return &CommentMongoRepository{
// 		DB:     db,
// 		ctx:    ctx,
// 		Cancel: cancel,
// 	}
// }

// func (repo *CommentMongoRepository) GetAll(post_id string) ([]*Comment, error) {
// 	result := &posts.Post{}
// 	objectId, errGettingObject := primitive.ObjectIDFromHex(id)
// 	if errGettingObject != nil {
// 		log.Println("Error getting object from id:", id)
// 		return nil, errGettingObject
// 	}
// 	filter := bson.M{"_id": objectId}
// 	err := repo.DB.FindOne(*repo.Ctx, filter).Decode(&result)
// 	if err == mongo.ErrNoDocuments {
// 		fmt.Println("record does not exist")
// 		return nil, ErrNoPost
// 	} else if err != nil {
// 		log.Fatal(err)
// 	}

// 	return result, nil

// 	return result, nil
// }

// func (repo *CommentMongoRepository) Add(post_id string, body string, user *user.User) (*Comment, error) {

// 	newComment := &Comment{
// 		ID:      strconv.Itoa(repo.lastID),
// 		Body:    body,
// 		Created: time.Now(),
// 		Author:  user,
// 		PostID:  post_id,
// 	}
// 	repo.lastID++

// 	repo.mu.Lock()
// 	repo.data = append(repo.data, newComment)
// 	repo.mu.Unlock()

// 	return newComment, nil
// }

// func (repo *CommentMongoRepository) Delete(post_id, comment_id string) (bool, error) {
// 	i := -1
// 	for idx, item := range repo.data {
// 		if item.ID == comment_id && item.PostID == post_id {
// 			i = idx
// 			break
// 		}
// 	}
// 	if i < 0 {
// 		return false, ErrNoComment
// 	}

// 	repo.mu.Lock()
// 	if i < len(repo.data)-1 {
// 		copy(repo.data[i:], repo.data[i+1:])
// 	}
// 	repo.data[len(repo.data)-1] = nil
// 	repo.data = repo.data[:len(repo.data)-1]
// 	repo.mu.Unlock()

// 	return true, nil
// }

// func (repo *CommentMongoRepository) DeleteAllByPost(post_id string) (bool, error) {
// 	newData := make([]*Comment, 0, len(repo.data))
// 	for _, item := range repo.data {
// 		if item.PostID != post_id {
// 			newData = append(newData, item)
// 		}
// 	}
// 	repo.data = newData
// 	return true, nil
// }
