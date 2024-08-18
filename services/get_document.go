package services

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/db"
)

func GetWholeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	if name == "" && id == "" {
		http.Error(w, "name or id parameter is required", http.StatusBadRequest)
		return
	}

	// NOT YET IMPLEMENTED
	// if id != "" && name != "" {
	// 	result, err := getDocumentByIDandName(w, id)
	// 	db.SendDocument(w, result, err, "id")
	if id != "" {
		result, err := getDocumentByID(w, id)
		db.SendDocument(w, result, err, "id")
	} else if name != "" {
		result, err := getDocumentByName(w, name)
		db.SendDocument(w, result, err, "name")
	}
}

func getDocumentByID(w http.ResponseWriter, id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// This sends an HTTP response with the error message, no return value
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return nil, err
		// is further error handling necessary here, context to close?
	}
	return db.CallFindOne(context.TODO(), bson.D{{"_id", objID}})
}

func getDocumentByName(w http.ResponseWriter, name string) (bson.M, error) {
	return db.CallFindOne(context.TODO(), bson.D{{"name", name}})
}
