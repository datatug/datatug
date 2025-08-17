package ui

import (
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

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

	content := sneatnav.NewPanelFromTextView(tui, textView)

	sneatv.DefaultBorder(textView.Box)
	textView.SetTitle(fileName)
	textView.SetTitleAlign(tview.AlignLeft)

	menu := newDataTugMainMenu(tui, settingsRootScreen)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
	return nil
}
