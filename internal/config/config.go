package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Postgres   Postgres   `yaml:"postgres"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type Postgres struct {
	DB_USER      string        `env:"DB_USER" env-required:"true"`
	DB_PASSWORD  string        `env:"DB_PASSWORD" env-required:"true"`
	DBName       string        `env:"DB_NAME" env-required:"true"`
	Host         string        `yaml:"host" env-default:"postgres"`
	Port         int           `yaml:"port" env-default:"5432"`
	SSLMode      string        `yaml:"sslmode" env-default:"disable"`
	MaxOpenConns int           `yaml:"max_open_conns" env-default:"25"`
	ConnLifeTime time.Duration `yaml:"conn_max_lifetime" env-default:"5m"`
}

type HTTPServer struct {
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: %v", err)
	}

	return &cfg
}
