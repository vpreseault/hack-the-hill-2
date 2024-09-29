package database

import "errors"

func (db *DB) StartTimer(sessionID string, timer Timer) error {
	dbModel, err := db.load()
	if err != nil {
		return err
	}

	if _, exists := dbModel.Sessions[sessionID]; !exists {
		return errors.New("session not found")
	}

	session := dbModel.Sessions[sessionID]
	session.Timer = timer
	dbModel.Sessions[sessionID] = session

	err = db.write(dbModel)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) StopTimer(sessionID string) error {
	return nil
}
