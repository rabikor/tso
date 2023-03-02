package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFile_Exists(t *testing.T) {
	assert.True(t, isFileExists("../.env"))
}

func TestIsFile_NotExists(t *testing.T) {
	assert.False(t, isFileExists("../.env.not-found"))
}

func TestParseEnv_FileNotExists(t *testing.T) {
	dotenvPath := "../.env.not-found"

	env := Env{}

	assert.ErrorContains(
		t,
		env.ParseEnv(dotenvPath),
		fmt.Sprintf("file [%s] with env vars was not found", dotenvPath),
	)
}

func TestParseEnv_FailedLoad(t *testing.T) {
	dotenvPath := "../go.mod"

	env := Env{}

	assert.Error(
		t,
		env.ParseEnv(dotenvPath),
	)
}

func TestParseEnv_FailedProcess(t *testing.T) {
	dotenvPath := "./_mocks/.env"

	env := Env{}

	assert.Error(
		t,
		env.ParseEnv(dotenvPath),
	)
}

func TestParseEnv_Success(t *testing.T) {
	dotenvPath := "../.env"

	env := Env{}

	assert.NoError(
		t,
		env.ParseEnv(dotenvPath),
	)
}
