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

func UpdateDocumentHandler(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("updateDocumentHandler")
			var updateData CharacterWithID
			err := json.NewDecoder(r.Body).Decode(&updateData);

			if(err != nil) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if(updateData.Id == primitive.ObjectID{000000000000000000000000}) {
				fmt.Println("Character ID was nil\n")

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
				return
			}

			fmt.Printf("Character ID: %v\n", updateData.Id)

			filter := bson.M{"_id": updateData.Id}
			update := bson.M{"$set": bson.M{
				"name": updateData.Name, 
				"role": updateData.Role,
				"stats": updateData.Stats,
				"hp": updateData.HP,
				"humanity": updateData.Humanity,
				"currentSkills": updateData.CurrentSkills,
				"currentEffects": updateData.CurrentEffects}}

			result, err := config.Collection.UpdateOne(context.TODO(), filter, update)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result.ModifiedCount == 0 {
				http.Error(w, "No document was updated. It may already have the desired value.", http.StatusNotModified)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Character with id %v updated successfully\n", updateData.Id)
		},
	)
}