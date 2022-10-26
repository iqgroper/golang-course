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
	Add(item *Comment) (uint32, error)
	GetAll(post_id uint) ([]*Comment, error)
	Delete(post_id, comment_id uint) (bool, error)
}
