package sqliteschema

import (
	"database/sql"
	"testing"
)

func TestNewSchemaProvider(t *testing.T) {
	if got := NewSchemaProvider(func() (*sql.DB, error) {
		return nil, nil
	}); got == nil {
		t.Errorf("NewSchemaProvider() = %v, want %v", got, nil)
	}
}
