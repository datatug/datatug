package dtviewers

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/datatug/datatug-core/pkg/schemer"
)

type DbContext interface {
	Driver() Driver
	Schema() schemer.SchemaProvider
	GetDB(ctx context.Context) (db dal.DB, err error)
}

type Driver struct {
	ID         string
	ShortTitle string
}

type DbContextBase struct {
	driver Driver
	schema schemer.SchemaProvider
	getDB  func(ctx context.Context) (dal.DB, error)
}

func (c *DbContextBase) Driver() Driver {
	return c.driver
}

func (c *DbContextBase) Schema() schemer.SchemaProvider {
	return c.schema
}

func (c *DbContextBase) GetDB(ctx context.Context) (db dal.DB, err error) {
	return c.getDB(ctx)
}
