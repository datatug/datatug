package sqliteschema

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/datatug/datatug-core/pkg/schemer"
)

func (s schemaProvider) GetForeignKeysReader(c context.Context, schema, table string) (r schemer.ForeignKeysReader, err error) {
	_ = schema
	var db *sql.DB
	if db, err = s.getSqliteDB(); err != nil {
		return nil, err
	}
	if table == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}
	sqlText := fmt.Sprintf("PRAGMA foreign_key_list('%s')", table)
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve columns: %w", err)
	}
	return &foreignKeysReader{rows: rows}, nil
}

func (s schemaProvider) GetForeignKeys(ctx context.Context, schema, table string) (foreignKeys []schemer.ForeignKey, err error) {
	var r schemer.ForeignKeysReader
	r, err = s.GetForeignKeysReader(ctx, schema, table)
	if err != nil {
		return
	}
	for {
		var fk schemer.ForeignKey
		if fk, err = r.NextForeignKey(); err != nil {
			return
		}

		foreignKeys = append(foreignKeys, fk)

		if err = ctx.Err(); err != nil {
			return
		}
	}
}

type foreignKeysReader struct {
	table string
	fk    schemer.ForeignKey
	rows  *sql.Rows
}

func (r *foreignKeysReader) NextForeignKey() (schemer.ForeignKey, error) {
	var err error
	for r.rows.Next() {
		if err = r.rows.Err(); err != nil {
			return schemer.ForeignKey{}, err
		}
	}
	var (
		id, seq                                    int
		table, from, to, onUpdate, onDelete, match string
	)
	if err = r.rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
		return schemer.ForeignKey{}, err
	}
	isNewFK := r.fk.To.Name != table
	if isNewFK {
		result := r.fk
		r.fk = schemer.ForeignKey{
			Name: "", // SQLite has not FK name
			From: schemer.FKAnchor{
				Name:    r.table,
				Columns: []string{from},
			},
			To: schemer.FKAnchor{
				Name:    table,
				Columns: []string{to},
			},
		}
		if result.To.Name != "" {
			return result, nil
		}
	}
	return r.fk, io.EOF
}

func (r *foreignKeysReader) Close() error {
	return r.rows.Close()
}
