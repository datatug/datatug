package firestoreviewer

import (
	"fmt"
	"os"

	"github.com/datatug/datatug-cli/apps/firestoreviewer/fsviewer"
)

func Run() {
	if _, err := fsviewer.NewApp(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
