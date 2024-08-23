package handlers

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"main/app/db"
)

func GetDocumentByID(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			objID, err := primitive.ObjectIDFromHex(id)

			if err != nil {
				// This sends an HTTP response with the error message, no return value
				http.Error(w, "Invalid id format", http.StatusBadRequest)
				// is further error handling necessary here, context to close?
			}
			result, err := db.CallFindOne(context.TODO(), bson.D{{"_id", objID}})
			db.SendDocument(w, result, err, "ID")
		},
	)
}

func GetDocumentByName(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := r.PathValue("name")
			result, err := db.CallFindOne(context.TODO(), bson.D{{"name", name}})
			db.SendDocument(w, result, err, "name")
		},
	)
}
