package sneatnav

import "github.com/rivo/tview"

type Screen interface {
	tview.Primitive
	Options() ScreenOptions
	Activate() error
	Close() error
	//GetTitle() string
	//Window() tview.Primitive
}

type ScreenOptions struct {
	fullScreen bool
}

func (o ScreenOptions) FullScreen() bool {
	return o.fullScreen
}

//func FullScreen() ScreenOptions {
//	return ScreenOptions{fullScreen: true}
//}
