package handlers

import (
	"net/http"

	"main/app/services"
)

func HPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		services.GetHPHandler(w, r)
	case http.MethodPost:
		services.UpdateHP(w, r)
	default:
		http.Error(w, "Method not allowed, GET and POST only", http.StatusMethodNotAllowed)
	}
}
