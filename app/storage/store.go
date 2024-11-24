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
	_store  *map[string]Entry
	_config *map[string]Entry
}

func NewStore() Store {
	store := make(map[string]Entry)
	config := make(map[string]Entry)
	return Store{
		_store:  &store,
		_config: &config,
	}
}

func (store *Store) Set(
	key string,
	value string,
	expiresInMs *int,
) {
	store.set(key, value, expiresInMs, store._store)
}

func (store *Store) SetConfig(
	key string,
	value string,
	expiresInMs *int,
) {
	store.set(key, value, expiresInMs, store._config)
}

func (store *Store) Get(
	key string,
) (Entry, bool) {
	return store.get(key, store._store)
}

func (store *Store) GetConfig(
	key string,
) (Entry, bool) {
	return store.get(key, store._config)
}

func (store *Store) set(
	key string,
	value string,
	expiresInMs *int,
	_store *map[string]Entry,
) {
	var expiresAt *time.Time
	expiresAt = nil

	if expiresInMs != nil {
		expires := time.Now().UTC().Add(time.Duration(*expiresInMs) * time.Millisecond)
		expiresAt = &expires
	}

	(*_store)[key] = Entry{
		Value:     value,
		Type:      "string",
		ExpiresAt: expiresAt,
	}
}

func (store *Store) get(
	key string,
	_store *map[string]Entry,
) (Entry, bool) {
	value, ok := (*_store)[key]

	if !ok {
		return Entry{}, false
	}

	if value.ExpiresAt != nil {
		now := time.Now().UTC()
		fmt.Println(now)
		fmt.Println(*value.ExpiresAt)
		if now.After(*value.ExpiresAt) {
			delete((*_store), key)
			return Entry{}, false
		}
		return value, true
	}

	return value, true
}
