package dtproject

import (
	"context"

	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
)

var _ ProjectContext = (*projectContext)(nil)

type projectContext struct {
	context.Context
	tui     *sneatnav.TUI
	config  *appconfig.ProjectConfig
	loader  datatug.ProjectsLoader
	project *datatug.Project
	projErr chan error
}

func (p projectContext) WatchProject() <-chan error {
	return p.projErr
}

func (p projectContext) TUI() *sneatnav.TUI {
	return p.tui
}

func (p projectContext) Config() *appconfig.ProjectConfig {
	return p.config
}

func (p projectContext) Project() *datatug.Project {
	return p.project
}

type ProjectContext interface {
	context.Context
	TUI() *sneatnav.TUI
	Config() *appconfig.ProjectConfig
	Project() *datatug.Project
	WatchProject() <-chan error
}

func NewProjectContext(
	tui *sneatnav.TUI,
	config *appconfig.ProjectConfig,
	loader datatug.ProjectsLoader,
) ProjectContext {
	if tui == nil {
		panic("tui cannot be nil")
	}
	if config == nil {
		panic("config cannot be nil")
	}
	if loader == nil {
		panic("loader cannot be nil")
	}

	ctx := &projectContext{
		tui:     tui,
		config:  config,
		loader:  loader,
		projErr: make(chan error, 1),
	}
	go func() {
		project, err := loader.LoadProject(ctx, config.ID)
		if project != nil {
			ctx.project = project
		}
		ctx.projErr <- err
	}()
	return ctx
}
