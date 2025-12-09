package dtproject

import (
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
