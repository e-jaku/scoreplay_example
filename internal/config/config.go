package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

const (
	CONFIG_NAME = "config"
	CONFIG_TYPE = "yaml"
	CONFIG_PATH = "."
)

type ServerConfig struct {
	Address         string        `mapstructure:"address" validate:"required,hostname_port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	RequestTimeout  time.Duration `mapstructure:"request_timeout"`
}

type DBConfig struct {
	Host       string `mapstructure:"host" validate:"required,hostname"`
	Port       int    `mapstructure:"port" validate:"required,numeric"`
	User       string `mapstructure:"user" validate:"required"`
	Password   string `mapstructure:"password" validate:"required"`
	DBName     string `mapstructure:"dbname" validate:"required"`
	SSLMode    string `mapstructure:"sslmode" validate:"oneof=disable require"`
	Migrations string `mapstructure:"migrations" validate:"required"`
}

type StorageConfig struct {
	Bucket      string `mapstructure:"bucket" validate:"required"`
	Endpoint    string `mapstructure:"endpoint" validate:"required"`
	AccessKeyID string `mapstructure:"access_key_id" validate:"required"`
	SecretKeyID string `mapstructure:"secret_key_id" validate:"required"`
}

func LoadConfig() (*ServerConfig, *DBConfig, *StorageConfig, error) {
	v := viper.New()
	v.SetConfigName(CONFIG_NAME)
	v.SetConfigType(CONFIG_TYPE)
	v.AddConfigPath(CONFIG_PATH)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, nil, nil, xerrors.Errorf("error reading config file: %w", err)
	}

	var (
		serverConfig  ServerConfig
		dbConfig      DBConfig
		storageConfig StorageConfig
	)

	if err := v.UnmarshalKey("server", &serverConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("unable to decode server config: %w", err)
	}
	if err := v.UnmarshalKey("db", &dbConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("unable to decode DB config: %w", err)
	}
	if err := v.UnmarshalKey("storage", &storageConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("unable to decode DB config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(serverConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("invalid server config: %v", err)
	}
	if err := validate.Struct(dbConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("invalid db config: %v", err)
	}
	if err := validate.Struct(storageConfig); err != nil {
		return nil, nil, nil, xerrors.Errorf("invalid storage config: %v", err)
	}

	return &serverConfig, &dbConfig, &storageConfig, nil
}
