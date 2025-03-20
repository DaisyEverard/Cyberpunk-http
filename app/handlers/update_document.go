package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/app/config"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateDocumentHandler(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var updateData Character
			err := json.NewDecoder(r.Body).Decode(&updateData);

			if(err != nil) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if(updateData.Id == nil) {
				result, err := config.Collection.InsertOne(context.TODO(), updateData)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusCreated)
			    fmt.Fprintf(w, "Character with id %s created successfully", result.InsertedID)
				return
			}

			fmt.Printf("%v", updateData.Id) 

			id := updateData.Id

			filter := bson.M{"_id": id}
			upsert := true
			options := options.UpdateOptions{Upsert: &upsert}

			result, err := config.Collection.UpdateOne(context.TODO(), filter, updateData, &options)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result.ModifiedCount == 0 {
				http.Error(w, "No document was updated. It may already have the desired value.", http.StatusNotModified)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Character with id %v updated successfully", id)
		},
	)
}