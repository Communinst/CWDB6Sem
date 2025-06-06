package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	HttpServer    `yaml:"http-server" env-required:"true"`
	Database      `yaml:"database"`
	CacheDatabase Database
	CloudDatabase Database
}

type HttpServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-required:"true"`
}

func LogConfig(cfg Config) {
	fmt.Printf("%s %s %s %s %s %s\n %s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.Address)
}
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: Database{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_NAME"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
		HttpServer: HttpServer{
			Address: fmt.Sprintf("0.0.0.0:%s", os.Getenv("GCNTNR_PORT")), // "8080"
			Timeout: 4 * time.Second,
		},
	}

	// Validate required fields
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD not set")
	}
	return cfg, nil
}
