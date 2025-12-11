package sqliteschema

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug-core/pkg/schemer"
)

var _ schemer.CollectionsProvider = (*collectionsProvider)(nil)

type collectionsProvider struct {
	db *sql.DB
}

func (v collectionsProvider) GetCollections(
	_ context.Context, parentKey *dal.Key,
) (reader schemer.CollectionsReader, err error) {
	r := new(collectionsReader)
	reader = r
	if parentKey == nil {
		r.rows, err = v.db.Query(`SELECT type, name, sql FROM sqlite_schema WHERE type in ('table', 'view') ORDER BY name`)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve list of SQLite tables and views: %w", err)
		}
	} else {
		switch strings.ToLower(parentKey.ID.(string)) {
		case "tables":
			r.collectionType = "table"
			r.rows, err = v.db.Query(`SELECT name, sql FROM sqlite_schema WHERE type  'table' ORDER BY name`)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve list of SQLite tables : %w", err)
			}
		case "views":
			r.collectionType = "view"
			r.rows, err = v.db.Query(`SELECT name, sql FROM sqlite_schema WHERE type  'view' ORDER BY name`)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve list of SQLite views : %w", err)
			}
		default:
			err = fmt.Errorf("parent key ID expected to be either 'tables' or 'views', got '%s'", parentKey.ID.(string))
			return
		}
	}
	return
}

//goland:noinspection SqlNoDataSourceInspection

var _ schemer.CollectionsReader = (*collectionsReader)(nil)

type collectionsReader struct {
	collectionType string
	rows           *sql.Rows
	i              int
}

func (s *collectionsReader) NextCollection() (c *datatug.CollectionInfo, err error) {
	s.i++
	if !s.rows.Next() {
		if err = s.rows.Err(); err != nil {
			err = fmt.Errorf("failed to retrieve db object row #%d: %w", s.i, err)
		}
		return
	}
	c = new(datatug.CollectionInfo)
	if s.collectionType == "" {
		err = s.rows.Scan(c.Name, c.SQL)
	} else {
		err = s.rows.Scan(c.DbType, c.Name, c.SQL)
	}
	if err != nil {
		return c, fmt.Errorf("failed to scan db row into CollectionInfo struct: %w", err)
	}
	return
}
