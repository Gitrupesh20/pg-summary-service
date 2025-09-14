package domain

import "time"

type RemoteDBDetails struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type LocalDBStats struct {
	MaxConnections        int
	MaxIdleConnections    int
	MaxConnectionLifeTime time.Duration
}
