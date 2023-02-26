package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_GetOffset(t *testing.T) {
	p := Pagination{Limit: 10, Page: 2}

	assert.Equal(t, 10, p.GetOffset())
}
