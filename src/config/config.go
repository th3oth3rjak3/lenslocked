package config

import "fmt"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostGresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "Jake",
		Password: "",
		Name:     "lenslocked_dev",
	}
}

func TestPostGresConfig() PostgresConfig {
	cfg := DefaultPostGresConfig()
	cfg.Name = "lenslocked_test"
	return cfg
}

func DefaultHashKeyConfig() string {
	return "energistically matrix cloud-centric experiences"
}

type AppConfig struct {
	Port          int
	Env           string
	HashKeyConfig string
	PostgresConfig
}

func DefaultConfig() AppConfig {
	return AppConfig{
		Port:           3000,
		Env:            "dev",
		HashKeyConfig:  DefaultHashKeyConfig(),
		PostgresConfig: DefaultPostGresConfig(),
	}
}

func (ac AppConfig) IsProd() bool {
	return ac.Env == "prod"
}

func (ac AppConfig) IsDev() bool {
	return !ac.IsProd()
}
