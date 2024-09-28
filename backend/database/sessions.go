package database

import (
	"errors"
)

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

func (db *DB) CreateSession(hostName string) (string, error) {
	dbModel, err := db.load()
	if err != nil {
		return "", err
	}
	
	sessionID := SessionID(generateID())
	dbModel.Sessions[sessionID] = SessionInfo{
		Users: []string{hostName},
		Timer: Timer{},
	}

	err = db.write(dbModel)
	if err != nil {
		return "", err
	}

	return string(sessionID), nil
}