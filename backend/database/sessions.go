package database

import "errors"

func (db *DB) GetSessionInfo(sessionId SessionID) (SessionInfo, error) {
	session := SessionInfo{}

	dbModel, err := db.load()
	if err != nil {
		return session, err
	}

	session, exists := dbModel.Sessions[sessionId]
	if !exists {
		return session, errors.New("session not found")
	}

	return session, nil
}