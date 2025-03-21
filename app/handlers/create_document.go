package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/app/config"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewDocumentHandler(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("newDocumentHandler")
			var updateData CharacterWithID
			err := json.NewDecoder(r.Body).Decode(&updateData);

			if(err != nil) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dataWihoutID := CharacterWithoutID{
				updateData.Name,
				updateData.Role,
				updateData.Stats,
				updateData.HP,
				updateData.Humanity,
				updateData.CurrentSkills,
				updateData.CurrentEffects,
			}

				result, err := config.Collection.InsertOne(context.TODO(), dataWihoutID)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				message := fmt.Sprintf("Character with id %s created successfully\n", result.InsertedID)

				w.WriteHeader(http.StatusCreated)
				idObject := struct{Message string; Id interface{}}{message, result.InsertedID}
				json.NewEncoder(w).Encode(idObject)
		},
	)
}