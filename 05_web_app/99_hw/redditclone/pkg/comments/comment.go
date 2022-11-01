package comments

import (
	"redditclone/pkg/user"
	"time"
)

type Comment struct {
	ID      uint       `json:"id,string"`
	Body    string     `json:"body"`
	Created time.Time  `json:"created"`
	Author  *user.User `json:"author"`
	PostID  uint       `json:"-"`
}

type CommentRepo interface {
	Add(post_id uint, body string, user *user.User) (*Comment, error)
	GetAll(post_id uint) ([]*Comment, error)
	Delete(post_id, comment_id uint) (bool, error)
	DeleteAllByPost(post_id uint) (bool, error)
}
