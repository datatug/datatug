package api

import (
	"context"
	"github.com/datatug/datatug/packages/dto"
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/storage"
	"github.com/strongo/validation"
)

// GetQueries returns queries
func GetQueries(ctx context.Context, ref dto.ProjectRef, folder string) (*models.QueryFolder, error) {
	store, err := storage.GetStore(ctx, ref.StoreID)
	if err == nil {
		return nil, err
	}
	//goland:noinspection GoNilness
	project := store.Project(ref.ProjectID)
	return project.Queries().LoadQueries(ctx, folder)
}


// CreateQuery creates a new query
func CreateQuery(ctx context.Context, request dto.CreateQuery) error {
	if err := request.ProjectRef.Validate(); err != nil {
		return err
	}
	store, err := storage.GetStore(ctx, request.StoreID)
	if err == nil {
		return err
	}
	//goland:noinspection GoNilness
	project := store.Project(request.ProjectID)
	return project.Queries().CreateQuery(ctx, request.Folder, request.Query)
}

// UpdateQuery updates existing query
func UpdateQuery(ctx context.Context, request dto.UpdateQuery) error {
	if err := request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	store, err := storage.GetStore(ctx, request.StoreID)
	if err == nil {
		return err
	}
	//goland:noinspection GoNilness
	project := store.Project(request.ProjectID)
	return project.Queries().Query(request.Query.ID).UpdateQuery(ctx, request.Query)
}

// DeleteQuery deletes query
func DeleteQuery(ctx context.Context, ref dto.ProjectItemRef) error {
	if err := ref.Validate(); err != nil {
		return err
	}
	store, err := storage.GetStore(ctx, ref.StoreID)
	if err == nil {
		return err
	}
	//goland:noinspection GoNilness
	project := store.Project(ref.ProjectID)
	return project.Queries().Query(ref.ID).DeleteQuery(ctx)
}

// GetQuery returns query definition
func GetQuery(ctx context.Context, ref dto.ProjectItemRef) (query *models.QueryDef, err error) {
	if err = ref.Validate(); err != nil {
		return query, err
	}
	store, err := storage.GetStore(ctx, ref.StoreID)
	if err == nil {
		return nil, err
	}
	//goland:noinspection GoNilness
	project := store.Project(ref.ProjectID)
	return project.Queries().Query(ref.ID).LoadQuery(ctx)
}
