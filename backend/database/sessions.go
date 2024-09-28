package database

import (
	"errors"
)

func (db *DB) GetSessionInfo(sessionID string) (SessionInfo, error) {
	session := SessionInfo{}

	dbModel, err := db.load()
	if err != nil {
		return session, err
	}

	session, exists := dbModel.Sessions[sessionID]
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
	
	sessionID := generateID()
	dbModel.Sessions[sessionID] = SessionInfo{
		Users: []string{hostName},
		Timer: Timer{},
	}

	err = db.write(dbModel)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (db *DB) AddUserToSession(sessionID string, userName string) error {
	dbModel, err := db.load()
	if err != nil {
		return err
	}

	sessionInfo, exists := dbModel.Sessions[sessionID]
	if !exists {
		return errors.New("session not found")
	}	

	sessionInfo.Users = append(sessionInfo.Users, userName)
	dbModel.Sessions[sessionID] = sessionInfo

	err = db.write(dbModel)
	if err != nil {
		return err
	}	

	return nil
}
