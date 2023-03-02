package router

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {
	v := newValidator()

	assert.Error(t, v.Validate(echo.Map{}))
}
