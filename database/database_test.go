package database

import (
	"testing"

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
