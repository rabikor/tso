package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
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

func NewEnv(dotenvPath string) (env Env, _ error) {
	if !isFileExists(dotenvPath) {
		return env, fmt.Errorf("file [%s] with env vars was not found", dotenvPath)
	}

	if err := godotenv.Load(dotenvPath); err != nil {
		return env, err
	}

	if err := envconfig.Process("", &env); err != nil {
		return env, err
	}

	return env, nil
}

func isFileExists(path string) bool {
	if f, err := os.Stat(path); os.IsNotExist(err) || f.IsDir() {
		return false
	}

	return true
}
