package storage

import (
	"context"
	"github.com/datatug/datatug/packages/models"
)

type DbCatalogsStore interface {
	Server() DbServerStore
	DbCatalog(id string) DbCatalogStore
}

type DbCatalogStore interface {
	Catalogs() DbCatalogsStore
	LoadDbCatalogSummary(ctx context.Context) (catalog *models.DbCatalogSummary, err error)
}
