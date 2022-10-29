package comments

import "redditclone/pkg/user"

type Comment struct {
	ID      uint
	Body    string
	Created string
	Author  *user.User
	PostID  uint
}

type CommentRepo interface {
	// GetByID(id uint) (*Comment, error)
	Add(post_id uint, body string, user *user.User) (*Comment, error)
	GetAll(post_id uint) ([]*Comment, error)
	Delete(post_id, comment_id uint) (bool, error)
	DeleteAllByPost(post_id uint) (bool, error)
}
