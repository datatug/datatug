package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
)

type WithProjectConfig interface {
	GetProjectConfig() *appconfig.ProjectConfig
}

//type withProjectConfig struct {
//	projectConfig *appconfig.ProjectConfig
//}
//
//func (v withProjectConfig) GetProjectConfig() *appconfig.ProjectConfig {
//	return v.projectConfig
//}

func goProjectScreen(tui *sneatnav.TUI, project *appconfig.ProjectConfig) {
	goEnvironmentsScreen(tui, project)
}
