package gcloudui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"google.golang.org/api/cloudresourcemanager/v3"
)

type GCloudContext struct {
	TUI      *sneatnav.TUI
	Projects []*cloudresourcemanager.Project
}

type CGProjectContext struct {
	GCloudContext
	Project *cloudresourcemanager.Project
}
