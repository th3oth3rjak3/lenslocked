package config

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "Jake",
		Password: "",
		Name:     "lenslocked_dev",
	}
}

func TestPostgresConfig() PostgresConfig {
	cfg := DefaultPostgresConfig()
	cfg.Name = "lenslocked_test"
	return cfg
}

func DefaultHashKeyConfig() string {
	return "energistically matrix cloud-centric experiences"
}

type AppConfig struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	HmacKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		Port:     3000,
		Env:      "dev",
		HmacKey:  DefaultHashKeyConfig(),
		Database: DefaultPostgresConfig(),
	}
}

func (ac AppConfig) IsProd() bool {
	return ac.Env == "prod"
}

func (ac AppConfig) IsDev() bool {
	return !ac.IsProd()
}

func LoadConfig(configRequired bool) AppConfig {
	file, err := os.Open("config.json")
	if err != nil && configRequired {
		panic(err)
	}
	if err != nil && !configRequired {
		fmt.Println("Using the default config")
		return DefaultConfig()
	}
	var ac AppConfig
	dec := json.NewDecoder(file)
	err = dec.Decode(&ac)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded from config file")
	return ac
}
