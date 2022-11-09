package handlers

import (
	"html/template"
	"net/http"

	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Info error parse float data type")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}
