package posts

import "redditclone/pkg/comments"

type Post struct {
	ID        uint
	Title     string
	Score     int
	VotesList []struct {
		User string
		Vote int
	}
	Votes            int
	Category         string
	CreatedDTTM      string
	Text             string
	URL              string
	Type             string
	UpvotePercentage int
	Views            uint
	Comments         []*comments.Comment
	Author           struct {
		Username string
		ID       uint
	}
}

type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(item *NewPost) (*Post, error)
	GetAllByCategory(category string) ([]*Post, error)
	GetByID(post_id uint) (*Post, error)
	UpVote(post_id uint, username string) (*Post, error)
	DownVote(post_id uint, username string) (*Post, error)
	Delete(post_id uint) (bool, error)
}
