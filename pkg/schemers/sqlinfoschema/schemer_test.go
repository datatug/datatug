package sqlinfoschema

import (
	"testing"

	"github.com/datatug/datatug-core/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestNewInformationSchema(t *testing.T) {
	var server = models.ServerReference{Driver: "sql"}
	v := NewInformationSchema(server)
	assert.EqualValues(t, server, v.server)
}
