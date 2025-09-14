package config

import (
	"fmt"
	"os"
	"pg-summary-service/internal/domain"
	"time"
)

type config struct {
	port                  string
	localDbUrl            string
	maxConnections        int
	maxIdleConnections    int
	maxConnectionLifeTime time.Duration
	externalDbUrl         string
	retries               int
	isDebug               bool
	logDir                string
	logFile               string
}

var conf config

// LoadConfig initializes the config.
// note: we can load these less sensitive values from config and sensitive values from os env, for now just hardcoded
// note: config is initialized once at startup; getters are safe for concurrent use
func LoadConfig() error {
	localDb := os.Getenv("LOCAL_DB_URL")
	if localDb == "" {
		fmt.Println("No external db url set, using default")
		localDb = "postgresql://neondb_owner:npg_xYGowIyT68DZ@ep-purple-mouse-a1cneoys-pooler.ap-southeast-1.aws.neon.tech/p_v1?sslmode=require"
	}

	externalDb := os.Getenv("EXTERNAL_API_URL")
	if externalDb == "" {
		fmt.Println("No external api url set, using default")
		//externalDb = "http://localhost:3000/api/summary" //"http://external-service.local/api/summary" // fallback default
	}

	port := getEnv("PORT", "8080") // default is fine

	// Assign to package-level conf
	conf = config{
		port:                  port,
		localDbUrl:            localDb,
		maxConnections:        5,
		maxIdleConnections:    2,
		maxConnectionLifeTime: 5 * time.Minute,
		externalDbUrl:         externalDb,
		retries:               3,
		isDebug:               false,
		logDir:                "./logs",
		logFile:               "server.log",
	}

	return nil
}

func GetLocalDbUrl() string {
	return conf.localDbUrl
}

func GetLogDir() string {
	return conf.logDir
}
func GetLogFile() string {
	return conf.logFile
}
func Debug() bool {
	return conf.isDebug
}

func GetExternalDbUrl() string {
	return conf.externalDbUrl
}

func GetDBStats() domain.LocalDBStats {
	return domain.LocalDBStats{
		MaxConnections:        conf.maxConnections,
		MaxIdleConnections:    conf.maxIdleConnections,
		MaxConnectionLifeTime: conf.maxConnectionLifeTime,
	}
}

func GetPort() string {
	return conf.port
}

func GetRetries() int {
	return conf.retries
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
