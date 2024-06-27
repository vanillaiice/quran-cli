package list

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// terminal is a wrapper around a term.Terminal
// that implements some redundant methods.
type terminal struct {
	term   *term.Terminal
	state  *term.State
	output *termenv.Output
	fd     int
}

// newTerminal returns a new terminal.
func newTerminal() (*terminal, error) {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}

	output := termenv.NewOutput(os.Stdout)
	output.SaveScreen()

	return &terminal{
		term:   term.NewTerminal(os.Stdin, ""),
		state:  state,
		output: output,
		fd:     fd,
	}, nil
}

// Write writes a byte slice to the terminal.
func (t *terminal) Write(b []byte) (int, error) {
	return t.term.Write(b)
}

// WriteString writes a string to the terminal.
func (t *terminal) WriteString(s string) (int, error) {
	return t.term.Write([]byte(s))
}

// WriteRepeat writes n copies of b to the terminal.
func (t *terminal) WriteRepeat(b []byte, n int) {
	t.Write(bytes.Repeat(b, n))
}

// WriteStringRepeat writes n copies of s to the terminal.
func (t *terminal) WriteStringRepeat(s string, n int) {
	t.WriteString(strings.Repeat(s, n))
}

// Restore restores the initial terminal state.
func (t *terminal) Restore() error {
	t.output.RestoreScreen()
	return term.Restore(t.fd, t.state)
}

// Size returns the current size of the terminal.
func (t *terminal) Size() (int, int) {
	w, h, err := term.GetSize(t.fd)
	if err != nil {
		w, h = 80, 24
	}
	return w, h
}

// ClearScreen clears the terminal screen.
func (t *terminal) ClearScreen() {
	t.output.ClearScreen()
}

// MoveCursor moves the cursor.
func (t *terminal) MoveCursor(x, y int) {
	t.output.MoveCursor(x, y)
}

// HideCursor hides the cursor.
func (t *terminal) HideCursor() {
	t.output.HideCursor()
}

// ShowCursor shows the cursor.
func (t *terminal) ShowCursor() {
	t.output.ShowCursor()
}

// SetTitle sets the terminal title.
func (t *terminal) SetTitle(title string) {
	t.output.SetWindowTitle(title)
}

// RemoveTitle removes the terminal title.
func (t *terminal) RemoveTitle() {
	t.output.Write([]byte("\033]0;\007"))
}

// AltScreen switches to the alternate screen.
func (t *terminal) AltScreen() {
	t.output.AltScreen()
}

// ExitAltScreen exits the alternate screen.
func (t *terminal) ExitAltScreen() {
	t.output.ExitAltScreen()
}

// Reset resets the terminal.
func (t *terminal) Reset() {
	t.output.Reset()
}

// Return writes a carriage return.
func (t *terminal) Return() {
	t.Write([]byte("\r"))
}

// NewLine writes a new line.
func (t *terminal) NewLine() {
	t.Write([]byte("\n"))
}

// Reverse reverses the foreground and background colors.
func (t *terminal) Reverse() {
	t.Write([]byte("\033[7m"))
}

// SetColor sets the terminal color.
func (t *terminal) SetColor(color int) {
	t.WriteString(fmt.Sprintf("\033[%dm", color))
}

// Faint makes the text faint.
func (t *terminal) Faint() {
	t.WriteString("\033[2m")
}

// Bold makes the text bold.
func (t *terminal) Bold() {
	t.WriteString("\033[1m")
}
