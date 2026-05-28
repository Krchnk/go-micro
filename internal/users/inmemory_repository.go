package users

import (
	"sort"
	"sync"
)

type InMemoryRepository struct {
	mu     sync.RWMutex
	nextID int64
	users  map[int64]User
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nextID: 1,
		users:  make(map[int64]User),
	}
}

func (r *InMemoryRepository) Create(name, email string) User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := User{
		ID:    r.nextID,
		Name:  name,
		Email: email,
	}
	r.users[user.ID] = user
	r.nextID++

	return user
}

func (r *InMemoryRepository) Update(id int64, name, email string) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return User{}, ErrUserNotFound
	}

	updated := User{
		ID:    id,
		Name:  name,
		Email: email,
	}
	r.users[id] = updated
	return updated, nil
}

func (r *InMemoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

func (r *InMemoryRepository) List() []User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]User, 0, len(r.users))
	for _, user := range r.users {
		result = append(result, user)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}
