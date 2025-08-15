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

	// Draw items horizontally within the header row.
	cursorX := x
	maxX := x + width
	for i, item := range b.items {
		if cursorX >= maxX {
			break
		}
		label := tview.Escape(item.GetTitle())
		var text string
		if i == focus {
			if b.HasFocus() {
				// Highlight selected item with background color when focused.
				text = "[black:yellow]" + label + "[-:-:-]"
			} else {
				// Focused item (last by default) remains normal when not focused.
				text = label
			}
		} else {
			// Dim for unfocused items.
			text = "[::d]" + label + "[-:-:-]"
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

// Focus is called when this primitive receives focus.
func (b *Breadcrumbs) Focus(delegate func(p tview.Primitive)) {
	// Always select the last item upon receiving focus.
	if len(b.items) > 0 {
		b.focusedIndex = len(b.items) - 1
	}
	b.Box.Focus(delegate)
}

// Blur is called when this primitive loses focus.
func (b *Breadcrumbs) Blur() {
	b.Box.Blur()
}

// InputHandler handles keyboard input for navigation and activation.
func (b *Breadcrumbs) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if len(b.items) == 0 {
			return
		}
		sel := b.focusedIndex
		if sel < 0 || sel >= len(b.items) {
			sel = len(b.items) - 1
		}
		switch event.Key() {
		case tcell.KeyLeft:
			if sel > 0 {
				b.focusedIndex = sel - 1
			}
		case tcell.KeyRight:
			if sel < len(b.items)-1 {
				b.focusedIndex = sel + 1
			}
		case tcell.KeyEnter:
			if sel >= 0 && sel < len(b.items) && b.items[sel] != nil {
				_ = b.items[sel].Action()
			}
		case tcell.KeyTab, tcell.KeyDown:
			// Pass focus to the next focusable primitive by letting the container handle it.
			// We indicate leaving by blurring ourselves and not consuming the event.
			// Some containers handle Tab/Down globally; ensure we don't trap it.
			setFocus(nil)
		case tcell.KeyBacktab, tcell.KeyUp:
			// Similarly, allow previous focus via container.
			setFocus(nil)
		}
	})
}

// MouseHandler handles selection via mouse and focusing on click.
func (b *Breadcrumbs) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return b.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		if action != tview.MouseLeftDown && action != tview.MouseLeftClick {
			return false, nil
		}
		x, y := event.Position()
		if !b.InInnerRect(x, y) {
			return false, nil
		}
		// Determine which item was clicked based on x coordinate.
		innerX, innerY, width, _ := b.GetInnerRect()
		_ = innerY
		cursorX := innerX
		maxX := innerX + width
		for i, item := range b.items {
			if cursorX >= maxX {
				break
			}
			label := tview.Escape(item.GetTitle())
			// Compute width based on rune count (no style tags in label after Escape).
			w := len([]rune(label))
			if x >= cursorX && x < cursorX+w {
				b.focusedIndex = i
				setFocus(b)
				return true, nil
			}
			cursorX += w
			// Separator width
			if i < len(b.items)-1 && cursorX < maxX {
				sep := b.separator
				sepW := len([]rune(sep))
				cursorX += sepW
			}
		}
		setFocus(b)
		return true, nil
	})
}
