package template

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

var tmpl = template.New("main").Option("missingkey=error")

func Render(s string, data any) (out string, err error) {
	t, err := tmpl.Parse(s)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", errors.Wrap(err, "execute template")
	}

	return buf.String(), nil
}
