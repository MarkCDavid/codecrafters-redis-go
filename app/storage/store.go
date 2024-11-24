package storage

import (
	"regexp"
	"strings"
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
	expiresAt *time.Time,
) {
	store.set(key, value, expiresAt, store._store)
}

func (store *Store) SetConfig(
	key string,
	value string,
	expiresAt *time.Time,
) {
	store.set(key, value, expiresAt, store._config)
}

func (store *Store) Get(
	key string,
) (Entry, bool) {
	return store.get(key, store._store)
}

func (store *Store) Keys(
	pattern string,
) ([]string, error) {
	return store.keys(pattern, store._store)
}

func (store *Store) KeysConfig(
	pattern string,
) ([]string, error) {
	return store.keys(pattern, store._config)
}

func (store *Store) GetConfig(
	key string,
) (Entry, bool) {
	return store.get(key, store._config)
}

func (store *Store) set(
	key string,
	value string,
	expiresAt *time.Time,
	_store *map[string]Entry,
) {
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
		if now.After(*value.ExpiresAt) {
			delete((*_store), key)
			return Entry{}, false
		}
		return value, true
	}

	return value, true
}

func (store *Store) keys(
	pattern string,
	_store *map[string]Entry,
) ([]string, error) {
	re, err := regexp.Compile(globToRegex(pattern))
	if err != nil {
		return nil, err
	}

	var matchingKeys []string

	for key := range *_store {
		if re.MatchString(key) {
			matchingKeys = append(matchingKeys, key)
		}
	}

	return matchingKeys, nil
}

func globToRegex(pattern string) string {
	regex := regexp.QuoteMeta(pattern)

	regex = strings.ReplaceAll(regex, `\*`, `.*`)
	regex = strings.ReplaceAll(regex, `\?`, `.`)
	regex = strings.ReplaceAll(regex, `\[`, `[`)
	regex = strings.ReplaceAll(regex, `\]`, `]`)
	regex = strings.ReplaceAll(regex, `\^`, `^`)
	regex = strings.ReplaceAll(regex, `\-`, `-`)

	return "^" + regex + "$"
}
