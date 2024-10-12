package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Mongo      MongoConfig
	Postgres   PostgresConfig
	TmdbApi    TmdbConfig
	PagesCount int `yaml:"pages_count"`
}

type MongoConfig struct {
	Uri        string `yaml:"uri"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type PostgresConfig struct {
	Uri string `yaml:"uri"`
}

type TmdbConfig struct {
	BaseUrl string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

func MustLoad() *Config {
	configPath := getConfigPath()
	if configPath == "" {
		log.Fatalln("Config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s ", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Unable to load config file: %v", err)
	}

	return &cfg
}

func getConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	return res
}
