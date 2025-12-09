package dtsettings

import (
	"github.com/datatug/datatug-cli/apps/datatug/datatugui"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatnav"
	"github.com/datatug/datatug-cli/pkg/sneatview/sneatv"
	"github.com/datatug/datatug-core/pkg/appconfig"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

func RegisterModule() {
	datatugui.RegisterMainMenuItem(datatugui.RootScreenSettings,
		datatugui.MainMenuItem{
			Text:     "Settings",
			Shortcut: 's',
			Action:   goSettingsScreen,
		})
}

func goSettingsScreen(tui *sneatnav.TUI, focusTo sneatnav.FocusTo) error {
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

	menu := datatugui.NewDataTugMainMenu(tui, datatugui.RootScreenSettings)
	tui.SetPanels(menu, content, sneatnav.WithFocusTo(sneatnav.FocusToMenu))
	if focusTo == sneatnav.FocusToContent {
		tui.App.SetFocus(content)
	}
	return nil
}
