package api

import (
	"context"

	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug-core/pkg/dto"
	"github.com/datatug/datatug-core/pkg/storage"
	"github.com/strongo/validation"
)

// AddDbServer adds db server to project
func AddDbServer(ctx context.Context, ref dto.ProjectRef, projDbServer datatug.ProjDbServer) error {
	store, err := storage.GetProjectStore(ctx, ref.StoreID, ref.ProjectID)
	if err != nil {
		return err
	}
	return store.SaveProjDbServer(ctx, &projDbServer)
}

// UpdateDbServer adds db server to project
//
//goland:noinspection GoUnusedExportedFunction
func UpdateDbServer(ctx context.Context, ref dto.ProjectRef, projDbServer datatug.ProjDbServer) error {
	store, err := storage.GetProjectStore(ctx, ref.StoreID, ref.ProjectID)
	if err != nil {
		return err
	}
	return store.SaveProjDbServer(ctx, &projDbServer)
}

// DeleteDbServer adds db server to project
func DeleteDbServer(ctx context.Context, ref dto.ProjectRef, dbServer datatug.ServerReference) (err error) {
	store, err := storage.GetProjectStore(ctx, ref.StoreID, ref.ProjectID)
	if err != nil {
		return err
	}
	return store.DeleteProjDbServer(ctx, dbServer.ID())
}

// GetDbServerSummary returns summary on DB server
func GetDbServerSummary(ctx context.Context, ref dto.ProjectRef, dbServer datatug.ServerReference) (*datatug.ProjDbServerSummary, error) {
	if err := dbServer.Validate(); err != nil {
		err = validation.NewBadRequestError(err)
		return nil, err
	}
	store, err := storage.GetProjectStore(ctx, ref.StoreID, ref.ProjectID)
	if err != nil {
		return nil, err
	}
	return store.LoadProjDbServerSummary(ctx, dbServer.ID())
}
