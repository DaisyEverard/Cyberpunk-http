package handlers

import (
	"net/http"
)

func WholeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetWholeDocumentHandler(w, r)
	default:
		http.Error(w, "Method not allowed, GET only", http.StatusMethodNotAllowed)
	}
}
