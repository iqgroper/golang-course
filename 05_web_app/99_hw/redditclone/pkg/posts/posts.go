package posts

import "redditclone/pkg/user"

type Post struct {
	ID               uint   `schema:"-"`
	Title            string `schema:"title,required"`
	Score            int
	Votes            int
	Category         string
	CreatedDTTM      string
	Text             string
	Type             string
	UpvotePercentage int
	Views            uint
	Author           *user.User
}

type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(item *Post) (uint, error)
	GetAllByCategory(category string) ([]*Post, error)
	GetByID(post_id uint) (*Post, error)
	UpVote(post_id uint) (int, error)
	DownVote(post_id uint) (*Post, error)
	Delete(post_id uint) (*Post, error)
}
