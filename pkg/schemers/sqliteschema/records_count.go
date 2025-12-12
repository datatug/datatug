package sqliteschema

import (
	"context"
	"database/sql"
	"fmt"
)

func (s schemaProvider) RecordsCount(_ context.Context, catalog, schema, object string) (count *int, err error) {
	_ = catalog
	var sqliteDB *sql.DB
	sqliteDB, err = s.getSqliteDB()
	if err != nil {
		return
	}
	query := fmt.Sprintf("SELECT COUNT(1) FROM [%s].[%s]", schema, object)
	var rows *sql.Rows
	rows, err = sqliteDB.Query(query)
	defer func() {
		err2 := rows.Close()
		if err == nil {
			err = err2
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("failed to get records count for %v.%v: %w", schema, object, err)
	}
	_ = rows.Next()
	count = new(int)
	return count, rows.Scan(count)
}
