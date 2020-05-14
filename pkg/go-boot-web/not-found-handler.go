package gobootweb

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"strings"
)

type customError struct {
	Status  int    `json:"status" xml:"status"`
	Method  string `json:"method" xml:"method"`
	URI     string `json:"uri" xml:"uri"`
	Message string `json:"message" xml:"message"`
}

var notFoundTemplate, _ = template.New("NotFound").Parse("<html><head><title>{{.Status}} - {{.Message}}</title></head><body><h3>{{.Status}} - {{.Message}}</h3><p>Method: {{.Method}}<br/>URI: <a href=\"{{.URI}}\">{{.URI}}</a></p></body></html>")

func myCustomHandler(w http.ResponseWriter, r *http.Request) {
	customError := customError{http.StatusNotFound, r.Method, r.RequestURI, "Not Found"}
	accepts := strings.Split(r.Header.Get("Accept"), ",")
	for _, accept := range accepts {
		cleanAccept := strings.TrimSpace(strings.ToLower(accept))
		if strings.Contains(cleanAccept, "json") {
			break
		}
		if strings.Contains(cleanAccept, "html") {
			w.Header().Set("Content-Type", accept)
			w.WriteHeader(http.StatusNotFound)
			notFoundTemplate.Execute(w, customError)
			return
		}
		if strings.Contains(cleanAccept, "xml") {
			w.Header().Set("Content-Type", accept)
			w.WriteHeader(http.StatusNotFound)
			xml.NewEncoder(w).Encode(customError)
			return
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(customError)
}
