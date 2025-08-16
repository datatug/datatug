package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

var _ tview.Primitive = (*settingsPanel)(nil)
var _ sneatnav.Cell = (*settingsPanel)(nil)

type settingsPanel struct {
	sneatnav.PanelBase
	textView *tview.TextView
}

func (p *settingsPanel) Draw(screen tcell.Screen) {
	p.textView.Draw(screen)
}

func goSettingsScreen(tui *sneatnav.TUI) error {
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
		PanelBase: sneatnav.NewPanelBaseFromTextView(tui, textView),
		textView:  textView,
	}
	defaultBorder(content.textView.Box)
	content.textView.SetTitle(fileName)
	content.textView.SetTitleAlign(tview.AlignLeft)

	menu := newDataTugMainMenu(tui, settingsRootScreen)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
	return nil
}
