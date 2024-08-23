package handlers

import (
	"net/http"
)

func HPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHPHandler(w, r)
	case http.MethodPost:
		UpdateHP(w, r)
	default:
		http.Error(w, "Method not allowed, GET and POST only", http.StatusMethodNotAllowed)
	}
}
