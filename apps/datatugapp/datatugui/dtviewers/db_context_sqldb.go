package dtviewers

import (
	"context"
	"database/sql"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo2sql"
	"github.com/datatug/datatug-core/pkg/schemer"
)

type SqlDBGetter func(ctx context.Context, driverName string) (sqlDB *sql.DB, err error)

type SqlDBContext struct {
	DbContextBase
	GetSqlDB SqlDBGetter
}

func NewSqlDBContext(driver Driver, getSqlDB SqlDBGetter, schema schemer.SchemaProvider) *SqlDBContext {
	return &SqlDBContext{
		DbContextBase: DbContextBase{
			driver: driver,
			schema: schema,
			getDB: func(_ context.Context) (dal.DB, error) { // ctx reserved for future use
				sqlLiteDB, err := getSqlDB(context.Background(), driver.ID)
				if err != nil {
					return nil, err
				}
				db := dalgo2sql.NewDatabase(sqlLiteDB, dal.NewSchema(nil, nil), dalgo2sql.Options{})
				return db, nil
			},
		},
		GetSqlDB: getSqlDB,
	}
}
