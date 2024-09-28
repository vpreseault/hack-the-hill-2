package database

import (
	"encoding/json"
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

func (db *DB) write(dbModel DBModel) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbModel)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}