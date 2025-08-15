package ui

import (
	"github.com/datatug/datatug-cli/apps/datatug/tapp"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

var _ tview.Primitive = (*settingsPanel)(nil)
var _ tapp.Cell = (*settingsPanel)(nil)

type settingsPanel struct {
	tapp.PanelBase
	textView *tview.TextView
}

func (p *settingsPanel) Draw(screen tcell.Screen) {
	p.textView.Draw(screen)
}

func newSettingsScreen(tui *tapp.TUI) tapp.Screen {
	return newDefaultLayout(tui, settingsRootScreen, func(tui *tapp.TUI) (tapp.Cell, error) {
		setting, _ := appconfig.GetSettings()

		content, _ := yaml.Marshal(setting)

		const fileName = " Config File: ~/.datatug.yaml"
		textView := tview.NewTextView().SetText(string(content))
		panel := &settingsPanel{
			PanelBase: tapp.NewPanelBase(tui, textView, textView.Box),
			textView:  textView,
		}
		defaultBorder(panel.textView.Box)
		panel.textView.SetTitle(fileName)
		panel.textView.SetTitleAlign(tview.AlignLeft)
		return panel, nil
	})
}
