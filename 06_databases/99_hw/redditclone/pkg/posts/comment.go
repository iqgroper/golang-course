package posts

import (
	"redditclone/pkg/user"
	"time"
)

type Comment struct {
	ID      string     `json:"id"`
	Body    string     `json:"body"`
	Created time.Time  `json:"created"`
	Author  *user.User `json:"author"`
	PostID  string     `json:"-"`
}

type CommentRepo interface {
	Add(post_id string, body string, user *user.User) (*Comment, error)
	GetAll(post_id string) ([]*Comment, error)
	Delete(post_id, comment_id string) (bool, error)
}
