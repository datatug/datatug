package api

import (
	"github.com/datatug/datatug/packages/dto"
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/storage"
	"github.com/strongo/validation"
)

func validateProjectInput(projectID string) (err error) {
	if projectID == "" {
		return validation.NewErrRequestIsMissingRequiredField("projectID")
	}
	return nil
}

// GetProjects return all projects
func GetProjects(storeID string) ([]models.ProjectBrief, error) {
	dal, err := storage.NewDatatugStore(storeID)
	if err != nil {
		return nil, err
	}
	return dal.GetProjects()
}

// GetProjectSummary returns project summary
func GetProjectSummary(ref dto.ProjectRef) (*models.ProjectSummary, error) {
	if ref.ProjectID == "" {
		return nil, validation.NewErrRequestIsMissingRequiredField("id")
	}
	store, err := storage.GetStore(ref.StoreID)
	if err == nil {
		return nil, err
	}
	//goland:noinspection GoNilness
	project := store.Project(ref.ProjectID)
	projectSummary, err := project.LoadProjectSummary()
	return &projectSummary, err
}

// GetProjectFull returns full project metadata
func GetProjectFull(ref dto.ProjectRef) (*models.DatatugProject, error) {
	store, err := storage.GetStore(ref.StoreID)
	if err == nil {
		return nil, err
	}
	//goland:noinspection GoNilness
	project := store.Project(ref.ProjectID)
	return project.LoadProject()
}
