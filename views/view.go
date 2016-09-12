package views

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// LayoutDir dir containing HTML layouts
var LayoutDir = "views/layouts"
var bootstrap *template.Template

// View definition
type View struct {
	Template *template.Template
	Layout   string
}

// ViewData wraps the externally-passed data, and provides global data too.
type ViewData struct {
	Flashes map[string]string
	Data    interface{}
}

var funcmap = template.FuncMap{
	"menulink": func(page_id string, url string, link_text string) template.HTML {
		menulink := "<li><a href={{.URL}}>{{.LinkText}}</a></li>"
		tmpl, err := template.New("menulink.gohtml").Parse(menulink)
		if err != nil {
			panic(err)
		}

		data := struct {
			URL      string
			LinkText string
		}{url, link_text}

		var doc bytes.Buffer
		if err := tmpl.Execute(&doc, data); err != nil {
			panic(err)
		}

		return template.HTML(doc.String())
	},
}

// NewView constructs a new View for the given layout name & view file.
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.New(layout).Funcs(funcmap).ParseFiles(files...)

	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Render returns the rendered html/template.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	vd := ViewData{
		Flashes: flashes(),
		Data:    data,
	}
	fmt.Println("Rendering layout:", v.Layout, "for URL", r.URL.String())
	return v.Template.ExecuteTemplate(w, v.Layout, vd)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "/*.gohtml")
	if err != nil {
		panic(err)
	}
	return files
}

var flashRotator int

func flashes() map[string]string {
	flashRotator++

	if flashRotator%3 == 0 {
		return map[string]string{
			"warning": "You are about to exceed your plan limts!",
		}
	}

	return map[string]string{}
}
