package posts

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	IdMongo          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ID               string             `json:"id"`
	Title            string             `json:"title"`
	Score            int                `json:"score"`
	VotesList        []VoteStruct       `json:"votes"`
	Category         string             `json:"category"`
	CreatedDTTM      time.Time          `json:"created"`
	Text             string             `json:"text,omitempty"`
	URL              string             `json:"url,omitempty"`
	Type             string             `json:"type"`
	UpvotePercentage int                `json:"upvotePercentage"`
	Views            uint               `json:"views"`
	Comments         []*Comment         `json:"comments"`
	Author           AuthorStruct       `json:"author"`
}

type VoteStruct struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type AuthorStruct struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}
type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(item *NewPost) (*Post, error)
	GetAllByCategory(category string) ([]*Post, error)
	GetByID(post_id string) (*Post, error)
	GetByUser(user_login string) ([]*Post, error)
	UpVote(post_id, username string) (*Post, error)
	DownVote(post_id, username string) (*Post, error)
	UnVote(post_id, username string) (*Post, error)
	Delete(post_id string) (bool, error)
}
