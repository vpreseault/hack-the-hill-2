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

	err = db.write(DBModel{
		Sessions: make(map[SessionID]SessionInfo),
	})

    return db, err
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

func (db *DB) load() (DBModel, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbModel := DBModel{}

	file, err := os.Open(db.path)
    if err != nil {
        fmt.Println("Error opening DB file:", err)
        return dbModel, err
    }
    defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dbModel)
	if err != nil {
        fmt.Println("Error decoding DB file:", err)
        return dbModel, err
    }

	return dbModel, nil

}