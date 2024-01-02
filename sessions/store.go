package sessions

import (
	"errors"
	"time"
)

type Store interface {
	//Store save a session
	Save(sid string, created time.Duration) error

	//Load loads a session from the store
	Load(sid string) Session

	Delete(sid string)

	Expire()
}

// ideally the implementation should be in a separate file

type MemStore struct {
	expiry time.Duration
	data   map[string]Session
}

func NewMemStore(exp time.Duration) *MemStore {
	m := &MemStore{
		expiry: exp,
		data:   make(map[string]Session),
	}
	go func() {
		for {
			time.Sleep(m.expiry)
			m.Expire()
		}
	}()
	return m
}

func (m *MemStore) Save(sid string, created time.Time) error {
	_, ok := m.data[sid]
	if !ok {
		m.data[sid] = NewMapSession(sid, created)
		return nil
	}
	return errors.New("session with given id already exists")
}

func (m *MemStore) Load(sid string) Session {
	v, ok := m.data[sid]
	if !ok {
		return nil
	}
	return v
}

func (m *MemStore) Delete(sid string) error {
	_, ok := m.data[sid]
	if !ok {
		return errors.New("no session with given id")
	}
	delete(m.data, sid)
	return nil
}

func (m *MemStore) Expire() {
	now := time.Now()
	for k := range m.data {
		if m.data[k].Expired(now) {
			delete(m.data, k)
		}
	}
}
