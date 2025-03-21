package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/app/config"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func DeleteDocumentHandler(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandleFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("DeleteDocumentHandler\n")

			var updateData map[string]primative.ObjectID
			err := json.NewDecoder(r.Body).Decode(&updateData);

			if(err != nil) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			result, err := config.Collection.DeleteOne(context.TODO(), bson.D{{"_id", updateData.id}})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			message := fmt.Sprintf("Character with id %s created deleted\n", result.DeletedCount)
			if result.DeletedCount == 0 {
				http.Error(w, "No document was deleted", http.StatusNotModified)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
	)
}



	

				

		