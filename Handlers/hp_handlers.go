package main

import (
	"net/http"
)

func getHPHandler(w http.ResponseWriter, r *http.Request, store CharacterStore) {
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	if (name == "" && id == "")  {
		http.Error(w, "name or id parameter is required", http.StatusBadRequest)
		return
	}

	// logic for only one query field at a time
	// implement multiple query fields at a time 
	if id != "" {
		result, err := getDocumentByID(w,store, id)
		sendDocument(w, result, err, "id")
	} else if name != "" {
		result, err := getDocumentByName(w,store, name)
		sendDocument(w, result, err, "name")
	}
}