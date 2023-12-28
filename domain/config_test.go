package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDotEnv(t *testing.T) {
	assert.True(t, loadDotEnv())
}
