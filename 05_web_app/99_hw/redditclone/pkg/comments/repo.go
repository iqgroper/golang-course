package comments

import (
	"errors"
	"fmt"
	"redditclone/pkg/user"
	"sync"
	"time"
)

var (
	ErrNoComment = errors.New("no such comment found")
)

type CommentMemoryRepository struct {
	lastID uint
	data   []*Comment
	mu     *sync.RWMutex
}

func NewMemoryRepo() *CommentMemoryRepository {
	return &CommentMemoryRepository{
		data:   make([]*Comment, 0, 10),
		lastID: 0,
		mu:     &sync.RWMutex{},
	}
}

func (repo *CommentMemoryRepository) GetAll(post_id uint) ([]*Comment, error) {
	result := make([]*Comment, 0, 10)
	for _, item := range repo.data {
		if item.PostID == post_id {
			result = append(result, item)
		}
	}

	return result, nil
}

func (repo *CommentMemoryRepository) Add(post_id uint, body string, user *user.User) (*Comment, error) {

	current_time := time.Now()
	newComment := &Comment{
		ID:      repo.lastID,
		Body:    body,
		Created: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", current_time.Year(), current_time.Month(), current_time.Day(), current_time.Hour(), current_time.Minute(), current_time.Second()),
		Author:  user,
		PostID:  post_id,
	}
	repo.lastID++

	repo.mu.Lock()
	repo.data = append(repo.data, newComment)
	repo.mu.Unlock()

	return newComment, nil
}

func (repo *CommentMemoryRepository) Delete(post_id, comment_id uint) (bool, error) {
	i := -1
	for idx, item := range repo.data {
		if item.ID == comment_id && item.PostID == post_id {
			i = idx
			break
		}
	}
	if i < 0 {
		return false, ErrNoComment
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

func (repo *CommentMemoryRepository) DeleteAllByPost(post_id uint) (bool, error) {
	newData := make([]*Comment, 0, len(repo.data))
	for _, item := range repo.data {
		if item.PostID != post_id {
			newData = append(newData, item)
		}
	}
	repo.data = newData
	return true, nil
}
