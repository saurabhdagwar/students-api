package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true" `
}

// Env-default for production:   env-default:"production"
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" `
	StoragePath string `yaml:"storage_path" env-required:"true"`
	Addr        string `yaml:"address" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configFlag := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *configFlag
	}

	if configPath == "" {
		configPath = filepath.Join("config", "local.yaml")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exits %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Can not read config file %s", err.Error())
	}
	return &cfg
}
