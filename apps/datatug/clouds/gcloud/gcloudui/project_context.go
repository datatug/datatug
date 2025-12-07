package gcloudui

import (
	"context"

	"github.com/datatug/datatug-cli/pkg/auth/gauth"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"google.golang.org/api/cloudresourcemanager/v3"
)

type GCloudContext struct {
	TUI      *sneatnav.TUI
	projects []*cloudresourcemanager.Project
}

func (v *GCloudContext) GetProjects() (project []*cloudresourcemanager.Project, err error) {
	if v.projects == nil {
		ctx := context.Background()
		if v.projects, err = gauth.GetGCloudProjects(ctx); err != nil {
			return
		}
	}
	return v.projects, nil
}

type CGProjectContext struct {
	*GCloudContext
	Project *cloudresourcemanager.Project
}
