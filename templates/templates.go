package templates

import (
	"embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed *.tmpl
var FS embed.FS

var templates *template.Template

func init() {
	templates = template.New("")

	if _, err := templates.ParseFS(FS, "*"); err != nil {
		panic(err)
	}

	fmt.Println(templates.DefinedTemplates())
}

func ExecuteOrderPlacedDE(w io.Writer, data any) error {
	return templates.ExecuteTemplate(w, "order_placed_de.tmpl", data)
}
