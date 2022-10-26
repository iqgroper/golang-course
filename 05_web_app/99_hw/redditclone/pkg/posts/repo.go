package posts

// WARNING! completly unsafe in multi-goroutine use, need add mutexes

type ItemMemoryRepository struct {
	lastID uint32
	data   []*Post
}

func NewMemoryRepo() *ItemMemoryRepository {
	return &ItemMemoryRepository{
		data: make([]*Post, 0, 10),
	}
}

func (repo *ItemMemoryRepository) GetAll() ([]*Post, error) {
	return repo.data, nil
}

func (repo *ItemMemoryRepository) GetByID(id uint32) (*Post, error) {
	for _, item := range repo.data {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, nil
}

func (repo *ItemMemoryRepository) Add(item *Post) (uint32, error) {
	repo.lastID++
	item.ID = repo.lastID
	repo.data = append(repo.data, item)
	return repo.lastID, nil
}

func (repo *ItemMemoryRepository) Update(newItem *Post) (bool, error) {
	for _, item := range repo.data {
		if item.ID != newItem.ID {
			continue
		}
		item.Title = newItem.Title
		item.Description = newItem.Description
		return true, nil
	}
	return false, nil
}

func (repo *ItemMemoryRepository) Delete(id uint32) (bool, error) {
	i := -1
	for idx, item := range repo.data {
		if item.ID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil // or the zero value of T
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}
