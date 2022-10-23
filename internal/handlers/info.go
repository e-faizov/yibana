package handlers

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("layout.html"))
}
func (m *MetricsHandlers) Info(w http.ResponseWriter, r *http.Request) {
	data := m.Store.GetAll()
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}
