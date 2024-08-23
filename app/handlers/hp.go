package handlers

import (
	"net/http"

	"main/app/handlers"
)

func HPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.GetHPHandler(w, r)
	case http.MethodPost:
		handlers.UpdateHP(w, r)
	default:
		http.Error(w, "Method not allowed, GET and POST only", http.StatusMethodNotAllowed)
	}
}
