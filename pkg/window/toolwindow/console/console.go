// Package console provides a graphical console for logging output while the app is running.
package console

import (
	"fmt"
	"os"

	"github.com/gucio321/HellSpawner/pkg/app/state"

	g "github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/window/toolwindow"
)

const (
	mainWindowW, mainWindowH = 600, 200
	lineW, lineH             = -1, -1
)

var _ toolwindow.ToolWindow = (*Console)(nil)

// Console represents a console
type Console struct {
	*toolwindow.ToolWindowBase
	outputText string
	fontFixed  *g.FontInfo
	logFile    *os.File
}

// Create creates a new console
func Create(fontFixed *g.FontInfo, x, y float32, logFile *os.File) *Console {
	result := &Console{
		fontFixed:      fontFixed,
		ToolWindowBase: toolwindow.New("Console", state.ToolWindowTypeConsole, x, y),
		logFile:        logFile,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result
}

// Build builds a console
func (c *Console) Build() {
	c.IsOpen(&c.Visible).
		Layout(c.GetLayout())
}

func (c *Console) GetLayout() g.Widget {
	return g.Style().SetFont(c.fontFixed).To(
		g.InputTextMultiline(&c.outputText).
			Size(lineW, lineH).
			Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoUndoRedo),
	)
}

// Write writes input on console, stdout and (if exists) to the log file
func (c *Console) Write(p []byte) (n int, err error) {
	msg := string(p) // convert message from byte slice into string

	c.outputText = msg + c.outputText // append message

	fmt.Print(msg) // print to terminal

	if c.logFile != nil {
		n, err = c.logFile.Write(p) // print to file
		if err != nil {
			return n, fmt.Errorf("error writing to log file: %w", err)
		} else if n != len(p) {
			return n, fmt.Errorf("invalid data written to log file")
		}
	}

	return len(p), nil
}
