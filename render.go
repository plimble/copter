package copter

import (
	"gopkg.in/unrolled/render.v1"
	"html"
	"html/template"
)

type Render struct {
	*render.Render
}

type Options render.Options

func New(options Options) *Render {
	options.Funcs = append(options.Funcs, template.FuncMap{"esp": esp})

	return &Render{render.New(render.Options(options))}
}

func esp(s string) string {
	return html.EscapeString(s)
}
