package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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