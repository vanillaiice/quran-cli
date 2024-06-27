package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/reflow/wordwrap"
	"github.com/vanillaiice/quran-cli/arabic"
	"github.com/vanillaiice/quran-cli/db"
	"github.com/vanillaiice/quran-cli/tui"
)

// Run runs the application.
func Run(s *db.Surah, l ...tui.Lang) (err error) {
	var lang tui.Lang

	if len(l) > 0 {
		lang = l[0]
	} else {
		lang = tui.Both
	}

	t, err := newTerminal()
	if err != nil {
		return
	}
	defer t.Restore()

	t.HideCursor()
	t.AltScreen()
	t.ClearScreen()
	t.MoveCursor(0, 0)
	t.SetTitle("Quran CLI")

	defer func() {
		t.Reset()
		t.ShowCursor()
		t.ExitAltScreen()
		t.ClearScreen()
		t.MoveCursor(0, 0)
		t.RemoveTitle()
	}()

	var currentLine, topLine int

	w, h := t.Size()

	printLines := func() {
		t.ClearScreen()
		t.Reset()
		t.MoveCursor(0, 0)

		var linesPrinted int

		for i := topLine; i < s.TotalVerses && linesPrinted < h-2; i++ {
			v := s.Verses[i]

			var s string

			switch lang {
			case tui.Arabic:
				s = fmt.Sprintf("%s. %s", arabic.ToArabic(v.Id), v.Text)
			case tui.Translation:
				s = fmt.Sprintf("%d. %s", v.Id, v.Translation)
			case tui.Both:
				fallthrough
			default:
				if v.Translation != "" {
					s = fmt.Sprintf("%d. %s\n%s. %s", v.Id, v.Translation, arabic.ToArabic(v.Id), v.Text)
				} else {
					s = fmt.Sprintf("%s. %s", arabic.ToArabic(v.Id), v.Text)
				}
			}

			wrapped := strings.Split(wordwrap.String(s, w-1), "\n")

			for j := 0; j < len(wrapped) && linesPrinted < h-2; j++ {
				wr := wrapped[j]

				if i == currentLine {
					t.Reverse()
					t.WriteString("|")
				}

				t.WriteString(wr + "\n")

				linesPrinted++

				t.Reset()
			}

			t.WriteString("\n")

			linesPrinted++
		}

		t.WriteStringRepeat("~\n", h-linesPrinted-1)

		t.Reverse()
		t.Bold()
		t.WriteStringRepeat(" ", w)
		b := []byte(fmt.Sprintf(" Verse %d/%d ", currentLine+1, s.TotalVerses))
		t.WriteStringRepeat("\b", len(b)-1)
		t.Write(b)
		t.WriteString(fmt.Sprintf("\r #%d %s (%s) - %s (%s) ", s.Id, s.Name, s.Transliteration, s.Translation, s.Type))
		// t.WriteString(" | ↑/k up • ↓/j down • q/esc exit • g/G top/bottom ")

		t.Reset()
	}

	printLines()

	done := make(chan error)

	go func() {
		up := func() {
			if currentLine > 0 {
				currentLine--
				if topLine > 0 {
					topLine--
				}
				printLines()
			}
		}

		down := func() {
			if currentLine < s.TotalVerses-1 {
				currentLine++
				if currentLine >= topLine+2 {
					topLine++
				}
				printLines()
			}
		}

		for {
			buf := make([]byte, 3)

			n, err := os.Stdin.Read(buf)
			if err != nil {
				done <- err
				return
			}

			switch n {
			case 1:
				switch buf[0] {
				case 'j':
					down()
				case 'k':
					up()
				case 'g':
					currentLine = 0
					topLine = 0
					printLines()
				case 'G':
					currentLine = s.TotalVerses - 1
					topLine = s.TotalVerses - 2
					printLines()
				case 'q', 27:
					done <- nil
				}
			case 3:
				if buf[0] == 27 && buf[1] == 91 {
					switch buf[2] {
					case 65:
						up()
					case 66:
						down()
					}
				}
			}
		}
	}()

	return <-done
}
