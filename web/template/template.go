package template

//go:generate go-bindata -pkg template -o template_gen.go files/

import (
	"encoding/json"
	"html/template"
	"path/filepath"
)

func Template() *template.Template {
	funcs := map[string]interface{}{
		"json": marshal,
	}

	dir, _ := AssetDir("files")
	tmpl := template.New("_")
	tmpl.Funcs(funcs)

	for _, name := range dir {
		path := filepath.Join("files", name)
		src := MustAsset(path)
		tmpl = template.Must(
			tmpl.New(name).Parse(string(src)),
		)
	}

	return tmpl
}

// marshal is a helper function to render data as JSON
// inside the tempalte.
func marshal(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}
