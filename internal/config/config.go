package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type ServerConfig struct {
	Address         string        `envconfig:"SERVER_ADDRESS" default:":8080"`
	ReadTimeout     time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"5s"`
	WriteTimeout    time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout     time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"15s"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"30s"`
	RequestTimeout  time.Duration `envconfig:"SERVER_REQUEST_TIMEOUT" default:"45s"`
}

type DBConfig struct {
	Host       string `envconfig:"DB_HOST" default:"localhost"`
	Port       int    `envconfig:"DB_PORT" default:"5432"`
	User       string `envconfig:"DB_USER" default:"postgres"`
	Password   string `envconfig:"DB_PASSWORD" default:"password"`
	DBName     string `envconfig:"DB_DBNAME" default:"example"`
	SSLMode    string `envconfig:"DB_SSLMODE" default:"disable"`
	Migrations string `envconfig:"DB_MIGRATIONS" default:"./migrations"`
}

type StorageConfig struct {
	Bucket      string `envconfig:"STORAGE_BUCKET" default:"example"`
	Endpoint    string `envconfig:"STORAGE_ENDPOINT" default:"localhost:9000"`
	AccessKeyID string `envconfig:"STORAGE_ACCESS_KEY_ID" default:"admin"`
	SecretKeyID string `envconfig:"STORAGE_SECRET_KEY_ID" default:"password"`
}

func LoadConfig() (*ServerConfig, *DBConfig, *StorageConfig, error) {
	var serverConfig ServerConfig
	var dbConfig DBConfig
	var storageConfig StorageConfig

	if err := envconfig.Process("", &serverConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("error loading server config: %w", err)
	}
	if err := envconfig.Process("", &dbConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("error loading DB config: %w", err)
	}
	if err := envconfig.Process("", &storageConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("error loading storage config: %w", err)
	}

	return &serverConfig, &dbConfig, &storageConfig, nil
}
