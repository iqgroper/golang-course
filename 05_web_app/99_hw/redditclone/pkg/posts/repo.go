package posts

import (
	"errors"
	"fmt"
	"redditclone/pkg/comments"
	"redditclone/pkg/user"
	"sync"
	"time"
)

var (
	ErrNoPost  = errors.New("no such post found")
	ErrNoCanDo = errors.New("method not allowed")
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

func (repo *PostMemoryRepository) GetByUser(user_login string) ([]*Post, error) {
	result := make([]*Post, 0, 10)
	for _, item := range repo.data {
		if item.Author.Username == user_login {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, ErrNoPost
	}

	return result, nil
}

func (repo *PostMemoryRepository) Add(item *NewPost) (*Post, error) {

	newPost := &Post{
		ID:               repo.lastID,
		Title:            item.Title,
		Score:            1,
		VotesList:        []VoteStruct{{item.Author.Login, 1}},
		Category:         item.Category,
		Comments:         make([]*comments.Comment, 0, 10),
		CreatedDTTM:      time.Now(),
		Text:             item.Text,
		URL:              item.URL,
		Type:             item.Type,
		UpvotePercentage: 100,
		Views:            0,
		Author:           AuthorStruct{item.Author.Login, item.Author.ID},
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
		if item.ID == id {
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

func percetageCount(votes []VoteStruct) int {
	votesCount := len(votes)
	if votesCount == 0 {
		return 0
	}
	positiveVotes := 0
	negativeVotes := 0
	for _, voter := range votes {
		if voter.Vote == 1 {
			positiveVotes++
		}

		if voter.Vote == -1 {
			negativeVotes++
		}
	}

	result := (positiveVotes * 100) / votesCount
	if result < 0 {
		fmt.Println("NEGATIVE VOTES")
	}

	return result
}

func (repo *PostMemoryRepository) UpVote(post_id uint, username string) (*Post, error) {
	for indexPost, item := range repo.data {
		if item.ID == post_id {
			for _, voter := range item.VotesList {
				if voter.User == username && voter.Vote == 1 {
					return nil, ErrNoCanDo
				}
				if voter.User == username && voter.Vote == -1 {
					repo.data[indexPost], _ = repo.UnVote(post_id, username)
					repo.data[indexPost], _ = repo.UpVote(post_id, username)
					return repo.data[indexPost], nil
				}
			}

			item.Score += 1

			item.VotesList = append(item.VotesList, VoteStruct{username, 1})
			item.UpvotePercentage = percetageCount(item.VotesList)

			return item, nil
		}
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepository) DownVote(post_id uint, username string) (*Post, error) {
	for indexPost, item := range repo.data {
		if item.ID == post_id {
			for _, voter := range item.VotesList {
				if voter.User == username && voter.Vote == -1 {
					return nil, ErrNoCanDo
				}
				if voter.User == username && voter.Vote == 1 {
					repo.data[indexPost], _ = repo.UnVote(post_id, username)
					repo.data[indexPost], _ = repo.DownVote(post_id, username)
					return repo.data[indexPost], nil
				}
			}
			item.Score -= 1

			item.VotesList = append(item.VotesList, VoteStruct{username, -1})
			item.UpvotePercentage = percetageCount(item.VotesList)

			return item, nil
		}
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepository) UnVote(post_id uint, username string) (*Post, error) {
	postIndexToRemove := -1
	voteIndexToRemove := -1
LOOP:
	for postIdx, item := range repo.data {
		if item.ID == post_id {
			postIndexToRemove = postIdx
			for idx, voter := range item.VotesList {
				if voter.User == username {
					voteIndexToRemove = idx

					item.Score -= voter.Vote

					break LOOP
				}
			}
		}
	}

	repo.mu.Lock()
	if voteIndexToRemove < len(repo.data[postIndexToRemove].VotesList)-1 {
		copy(repo.data[postIndexToRemove].VotesList[voteIndexToRemove:], repo.data[postIndexToRemove].VotesList[voteIndexToRemove+1:])
	}
	repo.data[postIndexToRemove].VotesList[len(repo.data[postIndexToRemove].VotesList)-1] = VoteStruct{}
	repo.data[postIndexToRemove].VotesList = repo.data[postIndexToRemove].VotesList[:len(repo.data[postIndexToRemove].VotesList)-1]
	repo.mu.Unlock()

	repo.data[postIndexToRemove].UpvotePercentage = percetageCount(repo.data[postIndexToRemove].VotesList)

	return repo.data[postIndexToRemove], nil
}
