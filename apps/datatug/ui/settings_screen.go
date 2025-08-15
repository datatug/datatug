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

func goSettingsScreen(tui *tapp.TUI) error {
	textView := tview.NewTextView()
	var settingsStr string
	setting, err := appconfig.GetSettings()
	if err != nil {
		settingsStr = err.Error()
	}

	if settingsStr == "" {
		data, err := yaml.Marshal(setting)
		if err != nil {
			settingsStr = err.Error()
		} else {
			settingsStr = string(data)
		}
	}

	const fileName = " Config File: ~/.datatug.yaml"
	textView.SetText(string(settingsStr))
	content := &settingsPanel{
		PanelBase: tapp.NewPanelBaseFromTextView(tui, textView),
		textView:  textView,
	}
	defaultBorder(content.textView.Box)
	content.textView.SetTitle(fileName)
	content.textView.SetTitleAlign(tview.AlignLeft)

	menu := newDataTugMainMenu(tui, settingsRootScreen)
	tui.SetPanels(menu, content)
	return nil
}
