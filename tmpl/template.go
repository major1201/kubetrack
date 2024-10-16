package tmpl

import (
	"bytes"
	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/pkg/errors"
)

// ExecuteTextTemplate a helper which quickly execute text template with data
func ExecuteTextTemplate(tmpl string, data any) (result string, err error) {
	engine, err := textTemplate.New("").Funcs(GetFuncMap().TextFuncMap()).Parse(tmpl)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = engine.Execute(buf, data)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	result = buf.String()
	return
}

// ExecuteHTMLTemplate a helper which quickly execute html template with data
func ExecuteHTMLTemplate(tmpl string, data any) (result string, err error) {
	engine, err := htmlTemplate.New("").Funcs(GetFuncMap().HTMLFuncMap()).Parse(tmpl)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = engine.Execute(buf, data)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	result = buf.String()
	return
}
