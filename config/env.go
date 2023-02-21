package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

var Env struct {
	Mode string `default:"debug" envconfig:"APP_ENV"`
	DB   struct {
		Host     string `default:"localhost" envconfig:"DB_HOST"`
		Port     uint   `default:"3306" envconfig:"DB_PORT"`
		User     string `default:"root" envconfig:"DB_USER"`
		Password string `required:"true" envconfig:"DB_PASSWORD"`
		Name     string `required:"true" envconfig:"DB_NAME"`
	}
	API struct {
		Request struct {
			Limit int `default:"20" envconfig:"API_REQUEST_LIMIT"`
			Page  int `default:"1" envconfig:"API_REQUEST_OFFSET"`
		}
	}
	Server struct {
		Port uint `default:"8000" envconfig:"SERVER_PORT"`
	}
}

func fileExists(path string) bool {
	if f, err := os.Stat(path); os.IsNotExist(err) || f.IsDir() {
		return false
	}

	return true
}

func init() {
	dotenvPath := "./.env"

	if dotenvPathEnv := os.Getenv("DOTENV_PATH"); dotenvPathEnv != "" {
		dotenvPath = dotenvPathEnv
	}

	if !fileExists(dotenvPath) {
		log.Panic(fmt.Sprintf("file [%s] with env vars was not found", dotenvPath))
		return
	}

	if err := godotenv.Load(dotenvPath); err != nil {
		log.Panic(err)
		return
	}

	if err := envconfig.Process("", &Env); err != nil {
		log.Panic(err)
		return
	}
}
