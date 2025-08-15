package breadcrumbs

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Breadcrumbs struct {
	*tview.Box
	items     []Breadcrumb
	separator string
}

func NewBreadcrumbs(options ...option) *Breadcrumbs {
	b := &Breadcrumbs{
		Box:       tview.NewBox(),
		items:     make([]Breadcrumb, 0, 8),
		separator: " > ",
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

	// Draw items horizontally within the header row.
	cursorX := x
	maxX := x + width
	for i, item := range b.items {
		if cursorX >= maxX {
			break
		}
		// Print the item title and advance the cursor by the number of cells used.
		_, printed := tview.Print(screen, item.GetTitle(), cursorX, y, maxX-cursorX, tview.AlignLeft, tcell.ColorYellow)
		cursorX += printed
		// Print separator between items if there is still space.
		if i < len(b.items)-1 && cursorX < maxX {
			_, sepPrinted := tview.Print(screen, b.separator, cursorX, y, maxX-cursorX, tview.AlignLeft, tcell.ColorGray)
			cursorX += sepPrinted
		}
	}
}
