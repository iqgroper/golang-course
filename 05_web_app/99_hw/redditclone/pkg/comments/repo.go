package comments

import (
	"errors"
	"sync"
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
		data: make([]*Comment, 0, 10),
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

func (repo *CommentMemoryRepository) Add(item *Comment) (uint, error) {
	item.ID = repo.lastID
	repo.lastID++

	repo.mu.Lock()
	repo.data = append(repo.data, item)
	repo.mu.Unlock()

	return repo.lastID, nil
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

// func (repo *CommentMemoryRepository) GetByID(id uint32) (*Post, error) {
// 	for _, item := range repo.data {
// 		if item.ID == id {
// 			return item, nil
// 		}
// 	}
// 	return nil, nil
// }
