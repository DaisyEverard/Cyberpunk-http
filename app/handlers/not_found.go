package handlers

import (
	"net/http"
)

func NotFound() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Sorry punk, can't find that one", http.StatusNotFound)
		},
	)
}