package database

import (
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

func CreateDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	file, err := os.Create(db.path)
    if err != nil {
        return nil, fmt.Errorf("error creating database: %v", err)
    }
    defer file.Close()
    return db, nil
}