package copter

import (
	"github.com/flosch/pongo2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Context map[string]interface{}

type Options struct {
	Directory     string
	Extensions    []string
	IsDevelopment bool
}

type Copter struct {
	options *Options
	Set     *pongo2.TemplateSet
}

func New(options *Options) *Copter {
	copter := &Copter{}
	if options.IsDevelopment {
		pongo2.DefaultSet.Debug = true
	}

	copter.options = options
	copter.Set = pongo2.DefaultSet
	copter.Set.SetBaseDirectory(options.Directory)
	copter.compileTemplates()
	return copter
}

func (copter *Copter) ExecW(name string, context map[string]interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	tpl, err := copter.Set.FromCache(name)
	if err != nil {
		panic(err)
	}
	if err := tpl.ExecuteWriter(context, w); err != nil {
		panic(err)
	}
}

func (copter *Copter) Exec(name string, context map[string]interface{}) string {
	tpl, err := copter.Set.FromCache(name)
	if err != nil {
		panic(err)
	}
	result, err := tpl.Execute(context)
	if err != nil {
		panic(err)
	}
	return result
}

func (copter *Copter) ExecByte(name string, context map[string]interface{}) []byte {
	tpl, err := copter.Set.FromCache(name)
	if err != nil {
		panic(err)
	}

	result, err := tpl.ExecuteBytes(context)
	if err != nil {
		panic(err)
	}

	return result
}

func (copter *Copter) compileTemplates() {
	dir := copter.options.Directory

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := ""
		if strings.Index(rel, ".") != -1 {
			ext = "." + strings.Join(strings.Split(rel, ".")[1:], ".")
		}

		for _, extension := range copter.options.Extensions {
			if ext == extension {
				name := (rel[0 : len(rel)-len(ext)])
				pongo2.Must(copter.Set.FromCache(filepath.ToSlash(name) + ext))
				break
			}
		}

		return nil
	})
}
