package template

import (
	"bytes"
	"text/template"

	"github.com/svrana/geniveev"
	"github.com/svrana/geniveev/builtins"
)

func Parse(name string, templateStr string, templateValues geniveev.TemplateValues) (string, error) {
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"Title": builtins.Title,
	}).Parse(templateStr)
	if err != nil {
		return "", err
	}
	out := &bytes.Buffer{}
	if err := tmpl.Execute(out, templateValues); err != nil {
		return "", err
	}
	return out.String(), nil
}
