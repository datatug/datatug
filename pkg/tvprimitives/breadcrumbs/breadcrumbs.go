package breadcrumbs

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Breadcrumbs struct {
	*tview.Box
	items        []Breadcrumb
	separator    string
	focusedIndex int // -1 means last item is considered focused
}

func NewBreadcrumbs(options ...option) *Breadcrumbs {
	b := &Breadcrumbs{
		Box:          tview.NewBox(),
		items:        make([]Breadcrumb, 0, 8),
		separator:    " > ",
		focusedIndex: -1,
	}
	for _, o := range options {
		o(b)
	}
	return b
}

func (b *Breadcrumbs) Push(item Breadcrumb) {
	b.items = append(b.items, item)
}

func (b *Breadcrumbs) Clear() {
	b.items = make([]Breadcrumb, 0, 8)
}

func (b *Breadcrumbs) Draw(screen tcell.Screen) {
	// Draw the base box (background, border, title, etc.).
	b.DrawForSubclass(screen, b)

	// Get the inner drawing area and keep text on a single header line.
	x, y, width, _ := b.GetInnerRect()
	if width <= 0 {
		return
	}

	// Determine which item is focused. Default to last.
	focus := b.focusedIndex
	if focus < 0 || focus >= len(b.items) {
		focus = len(b.items) - 1
	}

	// Draw items horizontally within the header row as buttons: "[ Title ]".
	cursorX := x
	maxX := x + width
	for i, item := range b.items {
		if cursorX >= maxX {
			break
		}
		label := item.GetTitle()
		text := tview.Escape(label) // ensure literal brackets, no tag parsing
		if i != focus {
			// Dim for unfocused items.
			text = "[::d]" + text + "[-:-:-]"
		}
		_, printed := tview.Print(screen, text, cursorX, y, maxX-cursorX, tview.AlignLeft, tcell.ColorYellow)
		cursorX += printed
		// Add a separator between items if there is still room.
		if i < len(b.items)-1 && cursorX < maxX {
			_, sp := tview.Print(screen, b.separator, cursorX, y, maxX-cursorX, tview.AlignLeft, tcell.ColorGray)
			cursorX += sp
		}
	}
}
