package storage

import (
	"fmt"
	"time"
)

type Entry struct {
	Value     string
	Type      string
	ExpiresAt *time.Time
}

type Store struct {
	_store *map[string]Entry
}

func NewStore() Store {
	store := make(map[string]Entry)
	return Store{
		_store: &store,
	}
}

func (store *Store) Set(
	key string,
	value string,
	expiresInMs *int,
) {
	var expiresAt *time.Time
	expiresAt = nil

	if expiresInMs != nil {
		expires := time.Now().UTC().Add(time.Duration(*expiresInMs) * time.Millisecond)
		expiresAt = &expires
	}

	(*store._store)[key] = Entry{
		Value:     value,
		Type:      "string",
		ExpiresAt: expiresAt,
	}
}

func (store *Store) Get(
	key string,
) (Entry, bool) {
	value, ok := (*store._store)[key]

	if !ok {
		return Entry{}, false
	}

	if value.ExpiresAt != nil {
		now := time.Now().UTC()
		fmt.Println(now)
		fmt.Println(*value.ExpiresAt)
		if now.After(*value.ExpiresAt) {
			delete((*store._store), key)
			return Entry{}, false
		}
		return value, true
	}

	return value, true
}
