package copter

import (
	"github.com/flosch/pongo2"
	"github.com/plimble/utils/pool"
	"net/http"
	"os"
	"path/filepath"
)

type Context map[string]interface{}

type Options struct {
	Directory     string
	Extensions    []string
	IsDevelopment bool
	PoolSize      int
}

type Copter struct {
	options *Options
	Set     *pongo2.TemplateSet
	pool    *pool.BufferPool
}

func New(options *Options) *Copter {
	copter := &Copter{}
	if options.IsDevelopment {
		pongo2.DefaultSet.Debug = true
	}

	if options.PoolSize == 0 {
		copter.pool = pool.NewBufferPool(100)
	} else {
		copter.pool = pool.NewBufferPool(options.PoolSize)
	}

	if len(options.Extensions) == 0 {
		options.Extensions = []string{".tpl", ".html"}
	}

	if options.Directory == "" {
		options.Directory = "./"
	}

	copter.options = options
	copter.Set = pongo2.DefaultSet
	copter.Set.SetBaseDirectory(options.Directory)
	copter.compileTemplates()
	return copter
}

func (copter *Copter) ExecW(name string, context map[string]interface{}, w http.ResponseWriter) {
	buf := copter.pool.Get()
	defer copter.pool.Put(buf)

	w.WriteHeader(http.StatusOK)

	tpl, err := copter.Set.FromCache(name)
	if err != nil {
		panic(err)
	}
	if err := tpl.ExecuteWriterUnbuffered(context, buf); err != nil {
		panic(err)
	}

	if _, err = w.Write(buf.Bytes()); err != nil {
		panic(err)
	}
}

func (copter *Copter) Exec(name string, context map[string]interface{}) string {
	buf := copter.pool.Get()
	defer copter.pool.Put(buf)

	tpl, err := copter.Set.FromCache(name)
	if err != nil {
		panic(err)
	}

	if err := tpl.ExecuteWriterUnbuffered(context, buf); err != nil {
		panic(err)
	}

	return buf.String()
}

func (copter *Copter) ExecByte(name string, context map[string]interface{}) []byte {
	buf := copter.pool.Get()
	defer copter.pool.Put(buf)

	tpl, err := copter.Set.FromCache(name)
	if err != nil {
		panic(err)
	}

	if err := tpl.ExecuteWriterUnbuffered(context, buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func (copter *Copter) compileTemplates() {
	dir := copter.options.Directory

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)

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
