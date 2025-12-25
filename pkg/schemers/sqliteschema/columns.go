package sqliteschema

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/datatug/datatug-core/pkg/schemer"
)

var _ schemer.ColumnsProvider = (*columnsProvider)(nil)

type columnsProvider struct {
	getSqliteDB func() (*sql.DB, error)
}

func (v columnsProvider) GetColumnsReader(_ context.Context, catalog string, filter schemer.ColumnsFilter) (schemer.ColumnsReader, error) {
	_ = catalog
	db, err := v.getSqliteDB()
	if err != nil {
		return nil, err
	}
	tableName := filter.CollectionRef.Name()
	if tableName == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}
	sqlText := fmt.Sprintf("PRAGMA table_info('%s')", tableName)
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve columns: %w", err)
	}
	return &columnsReader{rows: rows}, nil
}

func (v columnsProvider) GetColumns(ctx context.Context, catalog string, filter schemer.ColumnsFilter) ([]schemer.Column, error) {
	r, err := v.GetColumnsReader(ctx, catalog, filter)
	if err != nil {
		return nil, err
	}
	return schemer.ReadColumns(ctx, r)
}

var _ schemer.ColumnsReader = (*columnsReader)(nil)

type columnsReader struct {
	table string
	rows  *sql.Rows
}

func (s *columnsReader) NextColumn() (col schemer.Column, err error) {
	if !s.rows.Next() {
		if err = s.rows.Err(); err != nil {
			err = fmt.Errorf("failed to retrieve column row: %w", err)
			return
		}
		err = io.EOF
		return
	}
	//col = new(schemer.Column)
	col.TableName = s.table

	var pk, notNull bool
	var defaultVal sql.NullString

	if err = s.rows.Scan(
		&col.OrdinalPosition,
		&col.Name,
		&col.DbType,
		&notNull,
		&defaultVal,
		&pk,
	); err != nil {
		return col, fmt.Errorf("failed to scan INFORMATION_SCHEMA.COLUMNS row into TableColumn struct: %w", err)
	}
	col.IsNullable = !notNull
	return col, nil
}
