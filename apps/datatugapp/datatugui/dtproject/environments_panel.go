package dtproject

import (
	"fmt"

	"github.com/datatug/datatug/pkg/sneatview/sneatnav"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newEnvironmentsPanel(ctx ProjectContext) sneatnav.Panel {

	project := ctx.Project()
	environments, err := project.GetEnvironments(ctx)

	if err != nil {
		textView := tview.NewTextView()
		textView.SetText(err.Error())
		textView.SetTextColor(tcell.ColorRed)
		return sneatnav.NewPanel(ctx.TUI(), sneatnav.WithBox(textView, textView.Box))
	}

	list := tview.NewList()
	list.SetTitle(fmt.Sprintf("Environments (%d)", len(environments)))
	list.SetWrapAround(false)
	for _, environment := range environments {
		title := environment.Title
		if title == environment.ID {
			title = ""
		}
		list.AddItem(environment.ID, title, rune(environment.ID[0]), nil)
	}

	//list.AddItem("DEV", "Development", 'd', nil)
	//list.AddItem("QA", "Quality Assurance", 'q', nil)
	//list.AddItem("UAT", "User Acceptance Testing", 'u', nil)
	//list.AddItem("PROD", "Production", 'p', nil)

	return sneatnav.NewPanel(ctx.TUI(), sneatnav.WithBox(list, list.Box))
}
