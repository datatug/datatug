package breadcrumbs

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
)

// helper to read a full line from the screen
func readLine(screen tcell.Screen, y, width int) string {
	var b strings.Builder
	for x := 0; x < width; x++ {
		r, comb, _, _ := screen.GetContent(x, y)
		if r == 0 {
			// nothing drawn at this cell
			r = ' '
		}
		b.WriteRune(r)
		if len(comb) > 0 {
			for _, cr := range comb {
				b.WriteRune(cr)
			}
		}
	}
	return b.String()
}

func newSimScreen(t *testing.T, width, height int) tcell.Screen {
	t.Helper()
	s := tcell.NewSimulationScreen("UTF-8")
	if err := s.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	s.SetSize(width, height)
	return s
}

func TestNewBreadcrumbs_DefaultsAndOptions(t *testing.T) {
	bc := NewBreadcrumbs()
	if bc == nil {
		t.Fatalf("NewBreadcrumbs returned nil")
	}
	if bc.separator != " > " {
		t.Errorf("default separator = %q, want %q", bc.separator, " > ")
	}

	bc2 := NewBreadcrumbs(WithSeparator(" / "))
	if bc2.separator != " / " {
		t.Errorf("WithSeparator not applied, got %q", bc2.separator)
	}
}

func TestBreadcrumbs_PushAndClear(t *testing.T) {
	bc := NewBreadcrumbs()
	if len(bc.items) != 0 {
		t.Fatalf("initial items length = %d, want 0", len(bc.items))
	}
	bc.Push(NewBreadcrumb("A"))
	bc.Push(NewBreadcrumb("B"))
	if len(bc.items) != 2 {
		t.Fatalf("after Push, items length = %d, want 2", len(bc.items))
	}
	bc.Clear()
	if got := len(bc.items); got != 0 {
		t.Fatalf("after Clear, items length = %d, want 0", got)
	}
}

func TestBreadcrumbs_Draw_SingleLineNoBorder(t *testing.T) {
	width := 30
	height := 1
	s := newSimScreen(t, width, height)
	defer s.Fini()

	bc := NewBreadcrumbs()
	bc.Push(NewBreadcrumb("DataTug"))
	bc.Push(NewBreadcrumb("Projects"))
	bc.Push(NewBreadcrumb("Demo"))
	bc.SetRect(0, 0, width, height)

	bc.Draw(s)

	line := readLine(s, 0, width)
	expected := "DataTug > Projects > Demo"
	if got := strings.TrimRight(line, " "); !strings.HasPrefix(got, expected) {
		t.Fatalf("drawn line prefix mismatch:\n got: %q\nwant: %q", got, expected)
	}
}

func TestBreadcrumbs_Draw_RespectsInnerRectWithBorder(t *testing.T) {
	// Box with border: inner Y should be 1. Height must be at least 3 for border.
	width := 20
	height := 3
	s := newSimScreen(t, width, height)
	defer s.Fini()

	bc := NewBreadcrumbs()
	bc.SetBorder(true)
	bc.Push(NewBreadcrumb("A"))
	bc.Push(NewBreadcrumb("B"))
	bc.SetRect(0, 0, width, height)
	bc.Draw(s)

	// y=0 is border row, ensure no content there
	line0 := readLine(s, 0, width)
	if strings.Contains(line0, "A > B") {
		t.Fatalf("content drawn on border row (y=0): %q", strings.TrimRight(line0, " "))
	}
	// Validate content drawn within inner rect span on innerY.
	innerX, innerY, innerW, _ := bc.GetInnerRect()
	var b strings.Builder
	for x := innerX; x < innerX+innerW; x++ {
		r, comb, _, _ := s.GetContent(x, innerY)
		if r == 0 {
			r = ' '
		}
		b.WriteRune(r)
		for _, cr := range comb {
			b.WriteRune(cr)
		}
	}
	innerLine := strings.TrimRight(b.String(), " ")
	if !strings.HasPrefix(innerLine, "A > B") {
		t.Fatalf("content not drawn within inner rect (y=%d, x>=%d): %q", innerY, innerX, innerLine)
	}
}

func TestBreadcrumbs_Draw_TruncatesAtWidth(t *testing.T) {
	width := 10 // small width to force truncation
	height := 1
	s := newSimScreen(t, width, height)
	defer s.Fini()

	bc := NewBreadcrumbs(WithSeparator("/"))
	bc.Push(NewBreadcrumb("ABCDEFGHI"))
	bc.Push(NewBreadcrumb("XYZ"))
	bc.SetRect(0, 0, width, height)
	bc.Draw(s)

	line := readLine(s, 0, width)
	// Expected to start with the first title and possibly part of separator/title, but never exceed width
	trimmed := strings.TrimRight(line, " ")
	if len([]rune(line)) != width {
		t.Fatalf("line width %d != expected %d", len([]rune(line)), width)
	}
	if !strings.HasPrefix(trimmed, "ABCDEFGHI") && !strings.HasPrefix(trimmed, "ABCDEFGH") && !strings.HasPrefix(trimmed, "ABCDEFG") {
		t.Fatalf("unexpected truncation result: %q", trimmed)
	}
}
