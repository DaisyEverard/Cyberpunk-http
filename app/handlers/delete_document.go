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
type responseBody struct {
	Id string `bson:"id"`;
}

func DeleteDocumentHandler(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("DeleteDocumentHandler")

			var updateData responseBody
			err := json.NewDecoder(r.Body).Decode(&updateData);

			if(err != nil) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			fmt.Printf("Id raw : %+v\n", updateData.Id)
			Id, err := primitive.ObjectIDFromHex(updateData.Id)

			if err != nil {
				http.Error(w, "Invalid id format", http.StatusBadRequest)
			}

			fmt.Printf("Id hex : %+v\n", updateData.Id)
			filter := bson.D{primitive.E{Key: "_id", Value: Id}}

			result, err := config.Collection.DeleteOne(context.TODO(), filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result.DeletedCount == 0 {
				http.Error(w, "No document was deleted", http.StatusNotModified)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
	)
}



	

				

		