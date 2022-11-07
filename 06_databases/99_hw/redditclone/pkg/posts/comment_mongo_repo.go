package posts

import (
	"context"
	"fmt"
	"log"
	"redditclone/pkg/user"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentMongoRepository struct {
	DB     *mongo.Collection
	Ctx    *context.Context
	Cancel context.CancelFunc
	lastID int
}

func NewMongoRepo(db *mongo.Collection, ctx *context.Context, cancel context.CancelFunc) *CommentMongoRepository {
	return &CommentMongoRepository{
		DB:     db,
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func (repo *CommentMongoRepository) getPostByID(id string) (*Post, error) {
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

func (repo *CommentMongoRepository) GetAll(post_id string) ([]*Comment, error) {

	post, err := repo.getPostByID(post_id)
	if err != nil {
		log.Println("error getting comments", err.Error())
		return nil, err
	}

	return post.Comments, nil
}

func (repo *CommentMongoRepository) Add(post_id string, body string, user *user.User) (*Comment, error) {

	newComment := &Comment{
		ID:      strconv.Itoa(repo.lastID),
		Body:    body,
		Created: time.Now(),
		Author:  user,
		PostID:  post_id,
	}
	repo.lastID++

	post, err := repo.getPostByID(post_id)
	if err != nil {
		log.Println("error getting comments", err.Error())
		return nil, errors.Wrap(err, "adding comment")
	}

	post.Comments = append(post.Comments, newComment)

	id, _ := primitive.ObjectIDFromHex(post_id)
	res := repo.DB.FindOneAndReplace(*repo.Ctx, bson.M{"_id": id}, post)
	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return nil, ErrNoPost
	} else if res.Err() != nil {
		log.Fatal("FindOneAndReplace err", res.Err().Error())
	}

	return newComment, nil
}

func (repo *CommentMongoRepository) Delete(post_id, comment_id string) (bool, error) {
	post, err := repo.getPostByID(post_id)
	if err != nil {
		log.Println("error deleting a comment", err.Error())
		return false, err
	}
	comments := post.Comments

	removeID := -1
	for idx, comment := range comments {
		if comment.ID == comment_id {
			removeID = idx
		}
	}

	if removeID < len(comments)-1 {
		copy(comments[removeID:], comments[removeID+1:])
	}
	comments[len(comments)-1] = nil
	comments = comments[:len(comments)-1]

	post.Comments = comments

	id, _ := primitive.ObjectIDFromHex(post_id)
	res := repo.DB.FindOneAndReplace(*repo.Ctx, bson.M{"_id": id}, post)
	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		return false, ErrNoPost
	} else if res.Err() != nil {
		log.Fatal("FindOneAndReplace err", res.Err().Error())
	}

	return true, nil
}
