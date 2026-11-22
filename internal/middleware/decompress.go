package middleware

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strings"
)

func DecompressRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := strings.Join(r.Header.Values("Content-Encoding"), ",")
		encoding = strings.TrimSpace(encoding)

		if encoding == "" || strings.EqualFold(encoding, "identity") {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.EqualFold(encoding, "gzip") {
			writeEncodingError(w, http.StatusUnsupportedMediaType, "unsupported Content-Encoding")
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			writeEncodingError(w, http.StatusBadRequest, "invalid gzip header")
			return
		}
		defer reader.Close()

		r.Body = reader
		r.Header.Del("Content-Encoding")
		r.Header.Del("Content-Length")
		r.ContentLength = -1

		next.ServeHTTP(w, r)
	})
}

func writeEncodingError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}