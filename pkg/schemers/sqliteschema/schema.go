package sqliteschema

import (
	"context"
	"database/sql"

	"github.com/dal-go/dalgo/dal"
	"github.com/datatug/datatug-core/pkg/schemer"
)

// NewSchemaProvider creates a new SchemaProvider for MS SQL Server
func NewSchemaProvider(getSqliteDB func() (*sql.DB, error)) schemer.SchemaProvider {
	if getSqliteDB == nil {
		panic("getSqliteDB cannot be nil")
	}
	return schemaProvider{
		getSqliteDB: getSqliteDB,
	}
}

var _ schemer.SchemaProvider = (*schemaProvider)(nil)

type schemaProvider struct {
	columnsProvider
	constraintsProvider
	indexColumnsProvider
	indexesProvider
	getSqliteDB func() (*sql.DB, error)
}

func (s schemaProvider) GetCollections(_ context.Context, parent *dal.Key) (schemer.CollectionsReader, error) {
	if db, err := s.getSqliteDB(); err != nil {
		return nil, err
	} else {
		return getCollections(db, parent)
	}
}

func (schemaProvider) IsBulkProvider() bool {
	return true
}
