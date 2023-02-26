package database

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"treatment-scheme-organizer/config"
)

func TestOpen_FailedConnect(t *testing.T) {
	env := config.Env{}
	assert.NoError(t, env.ParseEnv("../.env"))

	env.DB.Host = "wrong_host"
	_, err := Open(&env)
	assert.Error(t, err)
}

func TestOpen_Success(t *testing.T) {
	env := config.Env{}
	assert.NoError(t, env.ParseEnv("../.env"))

	_, err := Open(&env)
	assert.NoError(t, err)
}

var sqlxDB *sqlx.DB

func setup() {
	env := config.Env{}
	_ = env.ParseEnv("../.env")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		env.DB.User,
		env.DB.Password,
		env.DB.Host,
		env.DB.Port,
		env.DB.Name,
	)

	sqlxDB, _ = sqlx.Open("mysql", dsn)
}
