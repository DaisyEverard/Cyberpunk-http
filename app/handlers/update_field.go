package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/app/config"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUpdateData(w http.ResponseWriter, r *http.Request, fieldName string) (update bson.D) {
	var updateData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received updateData: %+v", updateData)

	newValue, ok := updateData[fieldName]
	if !ok {
		http.Error(w, fieldName+" field is required and must be a number", http.StatusBadRequest)
		return
	}
	update = bson.D{{"$set", bson.M{fieldName: newValue}}}
	return update
}

func UpdateFieldByID(usersCollection *mongo.Collection, fieldName string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, "Invalid id format", http.StatusBadRequest)
				return
			}

			update := getUpdateData(w, r, fieldName)

			result, err := config.Collection.UpdateOne(context.TODO(), bson.D{{"_id", objID}}, update)

			if err == mongo.ErrNoDocuments || result.MatchedCount == 0 {
				http.Error(w, fmt.Sprintf("No document found with the id %s", objID), http.StatusNotFound)
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result.ModifiedCount == 0 {
				http.Error(w, "No document was updated. It may already have the desired value.", http.StatusNotModified)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, fieldName+" of character with id %s updated successfully", id)
		},
	)
}

func UpdateFieldByName(usersCollection *mongo.Collection, fieldName string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			update := getUpdateData(w, r, fieldName)
			name := r.PathValue("name")
			result, err := config.Collection.UpdateOne(context.TODO(), bson.D{{"name", name}}, update)
			if err == mongo.ErrNoDocuments || result.MatchedCount == 0 {
				http.Error(w, fmt.Sprintf("No document found with the name %s", name), http.StatusNotFound)
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result.ModifiedCount == 0 {
				http.Error(w, "No document was updated. It may already have the desired value.", http.StatusNotModified)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "HP of character with name %s updated successfully", name)
		},
	)
}
