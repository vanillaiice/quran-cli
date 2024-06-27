package tview

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vanillaiice/quran-cli/arabic"
	"github.com/vanillaiice/quran-cli/db"
	"github.com/vanillaiice/quran-cli/tui"
)

// Run runs the tview application.
func Run(surah *db.Surah, l ...tui.Lang) (err error) {
	var lang tui.Lang

	if len(l) > 0 {
		lang = l[0]
	} else {
		lang = tui.Both
	}

	_ = lang

	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetTitle(" The Holy Quran ")

	var i int
	drawFunc := func() {
		var s string

		for _, v := range surah.Verses {
			switch lang {
			case tui.Arabic:
				textView.SetTextAlign(tview.AlignRight)
				v.Text = replaceBrackets(v.Text)
				s += fmt.Sprintf(`["%d"]%s. %s[""]`+"\n", i, arabic.ToArabic(v.Id), v.Text)
				i++
			case tui.Translation:
				v.Translation = replaceBrackets(v.Translation)
				s += fmt.Sprintf(`["%d"]%d. %s[""]`+"\n", i, v.Id, v.Translation)
				i++
			case tui.Both:
				fallthrough
			default:
				v.Text = replaceBrackets(v.Text)
				v.Translation = replaceBrackets(v.Translation)
				s += fmt.Sprintf(`["%d"]%d. %s`+"\n", i, v.Id, v.Translation)
				s += fmt.Sprintf(`%s. %s[""]`+"\n", arabic.ToArabic(v.Id), v.Text)
				i++
			}

			s += "\n"
		}

		fmt.Fprint(textView, s)
	}

	drawFunc()

	helpText := " ↑/k up • ↓/j down • q/esc exit • g/G top/bottom"

	var sel int

	frame := tview.NewFrame(textView).
		AddText(helpText, false, tview.AlignLeft, tcell.ColorWhite).
		AddText(fmt.Sprintf(" #%d %s (%s) - %s (%s) | verse %d/%d", surah.Id, surah.Name, surah.Transliteration, surah.Translation, surah.Type, sel+1, surah.TotalVerses), false, tview.AlignLeft, tcell.ColorWhite)

	up := func() {
		if sel > 0 {
			sel--
		}
	}

	down := func() {
		if sel < i-1 {
			sel++
		}
	}

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'k':
				up()
			case 'j':
				down()
			case 'g':
				sel = 0
			case 'G':
				sel = i - 1
			case 'q':
				app.Stop()
			}
		case tcell.KeyEsc:
			app.Stop()
		case tcell.KeyUp:
			up()
		case tcell.KeyDown:
			down()
		}

		frame.Clear().
			AddText(helpText, false, tview.AlignLeft, tcell.ColorWhite).
			AddText(fmt.Sprintf(" #%d %s - %s (%s) | verse %d/%d", surah.Id, surah.Name, surah.Transliteration, surah.Type, sel+1, surah.TotalVerses), false, tview.AlignLeft, tcell.ColorWhite)

		textView.Highlight(fmt.Sprint(sel))

		textView.ScrollToHighlight()

		return event
	})

	textView.Highlight("0")

	textView.SetBorder(true).SetBorderAttributes(tcell.AttrDim)

	return app.SetRoot(frame, true).SetFocus(frame).Run()
}

// replaceBrackets replaces brackets with parentheses in a string.
func replaceBrackets(s string) string {
	s = strings.ReplaceAll(s, "[", "(")
	return strings.ReplaceAll(s, "]", ")")
}
