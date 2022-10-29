package posts

import (
	"errors"
	"redditclone/pkg/user"
	"sync"
	"time"
)

var (
	ErrNoPost = errors.New("no such post found")
)

type PostMemoryRepository struct {
	lastID uint
	data   []*Post
	mu     *sync.RWMutex
}

type NewPost struct {
	Type     string
	Title    string
	Text     string
	URL      string
	Category string
	Author   user.User `json:"-"`
}

func NewMemoryRepo() *PostMemoryRepository {
	return &PostMemoryRepository{
		lastID: 0,
		data:   make([]*Post, 0, 10),
		mu:     &sync.RWMutex{},
	}
}

func (repo *PostMemoryRepository) GetAll() ([]*Post, error) {
	return repo.data, nil
}

func (repo *PostMemoryRepository) GetByID(id uint) (*Post, error) {
	for _, item := range repo.data {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepository) Add(item *NewPost) (*Post, error) {

	newPost := &Post{
		ID:    repo.lastID,
		Title: item.Title,
		Score: 1,
		Votes: 1,
		VotesList: []struct {
			User string
			Vote uint
		}{{item.Author.Login, 1}},
		Category:         item.Category,
		CreatedDTTM:      time.Now().String(),
		Text:             item.Text,
		URL:              item.URL,
		Type:             item.Type,
		UpvotePercentage: 100,
		Views:            0,
		Author: struct {
			Username string
			ID       uint
		}{item.Author.Login, item.Author.ID},
	}
	repo.lastID++

	repo.mu.Lock()
	repo.data = append(repo.data, newPost)
	repo.mu.Unlock()

	return newPost, nil
}

func (repo *PostMemoryRepository) Delete(id uint) (bool, error) {
	i := -1
	for idx, item := range repo.data {
		if item.ID != id {
			i = idx
			break
		}
	}
	if i < 0 {
		return false, ErrNoPost
	}

	repo.mu.Lock()
	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]
	repo.mu.Unlock()

	return true, nil
}

func (repo *PostMemoryRepository) GetAllByCategory(category string) ([]*Post, error) {
	result := make([]*Post, 0, 10)
	for _, item := range repo.data {
		if item.Category == category {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, ErrNoPost
	}

	return result, nil
}

func (repo *PostMemoryRepository) UpVote(post_id uint) (*Post, error) {
	for _, item := range repo.data {
		if item.ID == post_id {
			item.Score += 1
			item.Votes += 1
			item.UpvotePercentage = item.Score / item.Votes
			// item.VotesList = append(item.VotesList, struct {
			// 	User string
			// 	Vote uint
			// }{, 1})
			return item, nil
		}
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepository) DownVote(post_id uint) (*Post, error) {
	for _, item := range repo.data {
		if item.ID == post_id {
			item.Score -= 1
			item.Votes += 1
			item.UpvotePercentage = item.Score / item.Votes
			return item, nil
		}
	}
	return nil, ErrNoPost
}
