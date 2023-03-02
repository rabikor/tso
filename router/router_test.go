package router

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	r := New()

	assert.IsType(t, &echo.Echo{}, r)
}
