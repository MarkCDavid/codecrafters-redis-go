package rdb

import (
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type RDB struct {
	Header    string
	Metadata  map[string]string
	Databases map[int]map[string]storage.Entry
}

func NewRDB() RDB {
	return RDB{
		Metadata:  make(map[string]string),
		Databases: make(map[int]map[string]storage.Entry),
	}
}

func (rdb *RDB) NewDB(index int) {
	rdb.Databases[index] = make(map[string]storage.Entry)
}
