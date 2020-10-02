package dictionary

import "errors"

type Entity struct {
	Id     uint16
	Field1 int64
	Field2 string
}

type Storage interface {
	Add(d Entity) error
	Get(id uint16) (Entity, error)
}

type mapStorage struct {
	maxId   uint16
	storage map[uint16]Entity
}

func NewMapStorage(maxId uint16) *mapStorage {
	return &mapStorage{
		maxId:   maxId,
		storage: make(map[uint16]Entity, maxId+1),
	}
}

func (s *mapStorage) Add(e Entity) error {
	if e.Id > s.maxId {
		return errors.New("the entity id is too big")
	}
	s.storage[e.Id] = e
	return nil
}

func (s *mapStorage) Get(id uint16) (Entity, error) {
	if id > s.maxId {
		return Entity{}, errors.New("the entity id is too big")
	}
	e, ok := s.storage[id]
	if !ok {
		return Entity{}, errors.New("the entity not found")
	}
	return e, nil
}

type arrayStorage struct {
	maxId   uint16
	storage []Entity
}

func NewArrayStorage(maxId uint16) *arrayStorage {
	return &arrayStorage{
		maxId:   maxId,
		storage: make([]Entity, maxId+1),
	}
}

func (s *arrayStorage) Add(e Entity) error {
	if e.Id > s.maxId {
		return errors.New("the entity id is too big")
	}
	s.storage[e.Id] = e
	return nil
}

func (s *arrayStorage) Get(id uint16) (Entity, error) {
	if id > s.maxId {
		return Entity{}, errors.New("the entity id is too big")
	}
	e := s.storage[id]
	if e.Id == 0 {
		return Entity{}, errors.New("the entity not found")
	}
	return e, nil
}
