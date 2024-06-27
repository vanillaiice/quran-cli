package md

import (
	"bytes"
	"text/template"

	"github.com/vanillaiice/quran-cli/db"
)

// tmpl is the template for the markdown file.
const tmpl = `# Surah {{ .Id }} {{ .Name }} - {{ .Transliteration }}
## {{ .Type }} - {{ .TotalVerses }} verses
	{{ range .Verses }}
		{{ .Id }}. {{ .Translation }}
		{{ .Text }}
		{{ end }}
	> Surah {{ .Id }}: {{ .Transliteration }}
`

// MakeTmpl generates the markdown template for the given surah.
func MakeTmpl(s *db.Surah) (t string, err error) {
	tmpl, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buf, s)
	if err != nil {
		return
	}

	return buf.String(), nil
}
