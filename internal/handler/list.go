package handler

import (
	"bytes"
	"html/template"
	"net/http"
)

const page = `
<!DOCTYPE html>
<html>
<body>
	<h1>Метрики</h1>
	<ul>
		{{range $name, $value := .Gauges}}
	<li>gauge: {{$name}} = {{$value}}</li>
{{end}}

{{range $name, $value := .Counters}}
	<li>counter: {{$name}} = {{$value}}</li>
{{end}}

	</ul>
</body>
</html>
`

type pageData struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (h *Handler) ListHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "Incorrect method", http.StatusMethodNotAllowed)
		return
	}

	gauses, counters := h.storage.GetAll()

	data := pageData{
		Gauges:   gauses,
		Counters: counters,
	}
	tmpl, err := template.New("metrics").Parse(page)
	if err != nil {
		http.Error(resp, "Template error", http.StatusInternalServerError)
		return
	}
	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, data)
	if err != nil {
		http.Error(resp, "Template error", http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err := buffer.WriteTo(resp); err != nil {
		return
	}
}
