package sqliteschema

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug-core/pkg/schemer"
)

func getCollections(db *sql.DB, parentKey *dal.Key) (reader schemer.CollectionsReader, err error) {
	r := new(collectionsReader)
	reader = r
	if parentKey == nil {
		//dal.From(dal.NewCollectionRef("sqlite_schema", "", nil)).
		//	Where(dal.Field("type").Equal("table", "view")).
		//	OrderBy(dal.Ascending(dal.Field("name")))
		r.rows, err = db.Query(`SELECT type, name, sql FROM sqlite_schema WHERE type in ('table', 'view') ORDER BY name`)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve list of SQLite tables and views: %w", err)
		}
	} else {
		switch strings.ToLower(parentKey.ID.(string)) {
		case "tables":
			r.collectionType = "table"
			r.rows, err = db.Query(`SELECT name, sql FROM sqlite_schema WHERE type  'table' ORDER BY name`)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve list of SQLite tables : %w", err)
			}
		case "views":
			r.collectionType = "view"
			r.rows, err = db.Query(`SELECT name, sql FROM sqlite_schema WHERE type  'view' ORDER BY name`)
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
	var name, dbType string
	if s.collectionType == "" {
		err = s.rows.Scan(&c.DbType, &name, &c.SQL)
	} else {
		err = s.rows.Scan(&name, &c.SQL)
	}
	if err != nil {
		return c, fmt.Errorf("failed to scan db row into CollectionInfo struct: %w", err)
	}
	var collectionType datatug.CollectionType
	switch strings.ToUpper(dbType) {
	case "TABLE":
		collectionType = datatug.CollectionTypeTable
	case "VIEW":
	}
	c.CollectionKey = datatug.NewCollectionKey(collectionType, name, "", "", nil)
	return
}
