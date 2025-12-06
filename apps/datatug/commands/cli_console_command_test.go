package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Execute(t *testing.T) {
	cmd := consoleCommand{}
	assert.Nil(t, cmd.Execute(nil))
}
