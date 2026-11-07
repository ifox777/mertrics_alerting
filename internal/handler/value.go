package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ValueHandler(resp http.ResponseWriter, req *http.Request) {
	value, err := h.metrics.Value(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
	if err != nil {
		http.Error(resp, "Not found", http.StatusNotFound)
		return
	}

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(resp, value)
}
