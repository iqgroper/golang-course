package posts

import (
	"redditclone/pkg/comments"
	"time"
)

type Post struct {
	ID               uint                `json:"id,string"`
	Title            string              `json:"title"`
	Score            int                 `json:"score"`
	VotesList        []VoteStruct        `json:"votes"`
	Category         string              `json:"category"`
	CreatedDTTM      time.Time           `json:"created"`
	Text             string              `json:"text,omitempty"`
	URL              string              `json:"url,omitempty"`
	Type             string              `json:"type"`
	UpvotePercentage int                 `json:"upvotePercentage"`
	Views            uint                `json:"views"`
	Comments         []*comments.Comment `json:"comments"`
	Author           AuthorStruct        `json:"author"`
}

type VoteStruct struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type AuthorStruct struct {
	Username string `json:"username"`
	ID       uint   `json:"id,string"`
}
type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(item *NewPost) (*Post, error)
	GetAllByCategory(category string) ([]*Post, error)
	GetByID(post_id uint) (*Post, error)
	GetByUser(user_login string) ([]*Post, error)
	UpVote(post_id uint, username string) (*Post, error)
	DownVote(post_id uint, username string) (*Post, error)
	UnVote(post_id uint, username string) (*Post, error)
	Delete(post_id uint) (bool, error)
}
