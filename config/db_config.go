package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultMaxIdleConnections = 10
const defaultMaxOpenConnections = 100

var dbConfig *DbConfig

type DbConfig struct {
	Host               string
	Port               string
	User               string
	Pass               string
	Name               string
	MaxIdleConnections int
	MaxOpenConnections int
}

func NewDbConfig() *DbConfig {
	if dbConfig == nil {
		dbConfig = initializeDbConfig()
	}
	return dbConfig
}

func initializeDbConfig() *DbConfig {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	maxIdleConnString := os.Getenv("DB_MAX_IDLE_CONNECTIONS")
	maxIdleConn, err := strconv.Atoi(maxIdleConnString)
	if err != nil {
		maxIdleConn = defaultMaxIdleConnections
	}

	maxOpenConnString := os.Getenv("DB_MAX_OPEN_CONNECTIONS")
	maxOpenConn, err := strconv.Atoi(maxOpenConnString)
	if err != nil {
		maxOpenConn = defaultMaxOpenConnections
	}

	return &DbConfig{
		Host:               host,
		Port:               port,
		User:               user,
		Pass:               pass,
		Name:               name,
		MaxIdleConnections: maxIdleConn,
		MaxOpenConnections: maxOpenConn,
	}
}
