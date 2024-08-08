package main

import (
	"net/http"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getDocumentHandler(w http.ResponseWriter, r *http.Request, store CharacterStore) {
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

func getDocumentByID(w http.ResponseWriter, store CharacterStore, id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// This sends an HTTP response with the error message, no return value
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return nil, err
		// is further error handling necessary here, context to close?
	}
	return store.FindOne(context.TODO(), bson.D{{"_id", objID}}) 
}

func getDocumentByName(w http.ResponseWriter, store CharacterStore, name string) (bson.M, error) {
	return store.FindOne(context.TODO(), bson.D{{"name", name}}) 
}

func sendDocument(w http.ResponseWriter, result bson.M, err error, queryType string) {
	if err == mongo.ErrNoDocuments {
		http.Error(w, fmt.Sprintf("No document found with that %s",queryType), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}