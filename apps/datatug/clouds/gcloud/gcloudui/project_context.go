package gcloudui

import (
	"context"
	"sync"

	"github.com/datatug/datatug-cli/pkg/auth/gauth"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"google.golang.org/api/cloudresourcemanager/v3"
)

type GCloudContext struct {
	TUI             *sneatnav.TUI
	loadingProjects sync.Mutex
	projects        []*cloudresourcemanager.Project
}

func (v *GCloudContext) GetProjects() (projects []*cloudresourcemanager.Project, err error) {
	if v.projects == nil {
		v.loadingProjects.Lock()
		defer v.loadingProjects.Unlock()
		if v.projects == nil {
			v.projects, err = gauth.GetGCloudProjects(context.Background())
		}
	}
	return v.projects, err
}

type CGProjectContext struct {
	*GCloudContext
	Project *cloudresourcemanager.Project
}
