package sneatv

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Tab represents a single tab.
type Tab struct {
	Title string
	tview.Primitive
}

// Tabs is a tab container implemented using tview.Pages.
type Tabs struct {
	*tview.Flex

	tabBar *tview.TextView
	pages  *tview.Pages

	tabs   []Tab
	active int
}

// NewTabs creates a new tab container.
func NewTabs() *Tabs {
	tabBar := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	pages := tview.NewPages()

	t := &Tabs{
		Flex:   tview.NewFlex().SetDirection(tview.FlexRow),
		tabBar: tabBar,
		pages:  pages,
	}

	tabBar.SetInputCapture(t.handleInput)

	t.
		AddItem(tabBar, 1, 0, false).
		AddItem(pages, 0, 1, true)

	return t
}

// AddTab adds a new tab.
func (t *Tabs) AddTab(tab Tab) {
	index := len(t.tabs)
	t.tabs = append(t.tabs, tab)

	t.pages.AddPage(
		tab.Title,
		tab.Primitive,
		true,
		index == 0,
	)

	if index == 0 {
		t.active = 0
	}

	t.renderTabs()
}

// SwitchTo switches to a tab by index.
func (t *Tabs) SwitchTo(index int) {
	if index < 0 || index >= len(t.tabs) {
		return
	}

	t.active = index
	t.pages.SwitchToPage(t.tabs[index].Title)
	t.renderTabs()
}

// renderTabs redraws the tab bar.
func (t *Tabs) renderTabs() {
	t.tabBar.Clear()

	for i, tab := range t.tabs {
		if i == t.active {
			_, _ = fmt.Fprintf(
				t.tabBar,
				`["%d"][black:white] %s [-:-][""] `,
				i,
				tab.Title,
			)
		} else {
			_, _ = fmt.Fprintf(
				t.tabBar,
				`["%d"] %s [""] `,
				i,
				tab.Title,
			)
		}
	}
}

// handleInput handles keyboard navigation.
func (t *Tabs) handleInput(ev *tcell.EventKey) *tcell.EventKey {
	switch ev.Key() {
	case tcell.KeyRight:
		t.SwitchTo((t.active + 1) % len(t.tabs))
		return nil
	case tcell.KeyLeft:
		t.SwitchTo((t.active - 1 + len(t.tabs)) % len(t.tabs))
		return nil
	default:
		if ev.Modifiers() == tcell.ModAlt {
			if ev.Rune() >= '1' && ev.Rune() <= '9' {
				t.SwitchTo(int(ev.Rune() - '1'))
				return nil
			}
		}
		return ev
	}
}
