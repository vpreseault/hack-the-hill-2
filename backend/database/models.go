package database

type SessionInfo struct {
	Users []string `json:"users"`
	Timer Timer    `json:"timer"`
}

type Timer struct {
	StartTime         string `json:"start_time"`
	DurationInMinutes int64  `json:"duration_in_minutes"`
	Type              string `json:"type"`
}

type SessionID string

type DBModel struct {
	Sessions map[SessionID]SessionInfo `json:"sessions"`
}