package posts

import "redditclone/pkg/user"

type Post struct {
	ID               uint32 `schema:"-"`
	Title            string `schema:"title,required"`
	Score            int
	Category         string
	CreatedDTTM      string
	Text             string
	Type             string
	UpvotePercentage float32
	Views            uint
	Author           *user.User
}

type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(item *Post) (uint, error)
	GetAllByCategory(category string) ([]*Post, error)
	GetByID(post_id uint) (*Post, error)
	UpVote(post_id uint) (int, error)
	DownVote(post_id uint) (int, error)
	Delete(post_id uint) (bool, error)
}
