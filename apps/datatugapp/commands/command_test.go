package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandIsNotNil(t *testing.T) {
	cmd := DatatugCommand()
	assert.NotNil(t, cmd)
}
