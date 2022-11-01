package handlers

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

const templateString = `<ul>
    {{range .}}{{if eq .MType "gauge"}}<li>{{.ID}} - {{.Value}}</li>{{else}}<li>{{.ID}} - {{.Delta}}</li>{{end}}{{end}}
</ul>`

func init() {
	tmpl = template.Must(template.New("info").Parse(templateString))
}
func (m *MetricsHandlers) Info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/html")
	data, err := m.Store.GetAll(ctx)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}
