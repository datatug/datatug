package filestore

import (
	"context"
	"fmt"
	"github.com/datatug/datatug/packages/dto"
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/storage"
)

// NewStore create a storage for multiple projects by their dir paths
func NewStore(id string, pathsByID map[string]string) (fsStore storage.Store, err error) {
	return newStore(id, pathsByID), nil
}

var _ storage.Store = (*FsStore)(nil)

// FsStore provides implementation of file system storage
type FsStore struct {
	id               string
	pathByID         map[string]string
	fileSystemLoader // TODO: To be deleted
	//storeSaver       // TODO: To be deleted
}

func (store FsStore) CreateProject(_ context.Context, _ dto.CreateProjectRequest) (*models.ProjectSummary, error) {
	panic("not implemented")
}

func (store FsStore) GetProjectStore(id string) storage.ProjectStore {
	path := store.pathByID[id]
	return newFsProjectStore(id, path)
}

func (store FsStore) DeleteProject(_ context.Context, id string) error {
	return fmt.Errorf("not implemented yet, id=%s", id)
}

// GetProjects returns list of projects
func (store FsStore) GetProjects(context.Context) (projectBriefs []models.ProjectBrief, err error) {
	projectBriefs = make([]models.ProjectBrief, len(store.pathByID))
	var i int
	for id, path := range store.pathByID {
		projectBriefs[i] = models.ProjectBrief{}
		projectBriefs[i].ID = id
		projFile, err := LoadProjectFile(path)
		if err != nil {
			return projectBriefs, fmt.Errorf("failed to load project file: %w", err)
		}
		projectBriefs[i].Title = projFile.Title
		projectBriefs[i].Repository = projFile.Repository
		i++
	}
	return
}

// newStore creates an instance of storage that implements storage.Store
func newStore(id string, pathByID map[string]string) *FsStore {
	return &FsStore{
		id:       id,
		pathByID: pathByID,
	}
}

// NewSingleProjectStore creates an instance of storage that implements storage.Store for a single project
func NewSingleProjectStore(projectPath, projectID string) (storeInterface *FsStore, projID string) {
	if projectID == "" {
		projID = storage.SingleProjectID
	} else {
		projID = projectID
	}
	const storeID = "single_project_file_store"
	storeInterface = newStore(storeID, map[string]string{projID: projectPath})
	return
}
