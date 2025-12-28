package store

import (
	"errors"
	"sync"
)

type Item struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Backend defines the operations required by the handlers.
type Backend interface {
	Create(Item) Item
	List() []Item
	Get(int64) (Item, bool)
	Update(int64, Item) (Item, error)
	Delete(int64) bool
}

type Store struct {
	mu     sync.RWMutex
	items  map[int64]Item
	nextID int64
}

func New() *Store {
	return &Store{items: make(map[int64]Item), nextID: 1}
}

func (s *Store) Create(it Item) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	it.ID = s.nextID
	s.nextID++
	s.items[it.ID] = it
	return it
}

func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]Item, 0, len(s.items))
	for _, v := range s.items {
		res = append(res, v)
	}
	return res
}

func (s *Store) Get(id int64) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	return it, ok
}

func (s *Store) Update(id int64, in Item) (Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	if !ok {
		return Item{}, errors.New("not found")
	}
	in.ID = id
	s.items[id] = in
	return in, nil
}

func (s *Store) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	if !ok {
		return false
	}
	delete(s.items, id)
	return true
}
