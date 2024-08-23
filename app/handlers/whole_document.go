package handlers

import (
	"main/app/services"
	"net/http"
)

func WholeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		services.GetWholeDocumentHandler(w, r)
	default:
		http.Error(w, "Method not allowed, GET only", http.StatusMethodNotAllowed)
	}
}
