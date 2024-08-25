package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/db"
)

func GetJSONObjectByID(usersCollection *mongo.Collection, fieldName string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.D{
				{fieldName, 1},
			}
			findOptions := options.FindOne().SetProjection(projection)

			id := r.PathValue("id")
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, "Invalid id format", http.StatusBadRequest)
				return
			}
			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"_id", objID}}, findOptions)
			if err != nil {
				http.Error(w, "Doucment with ID "+id+" not found", http.StatusNotFound)
				return
			}

			fieldValue, ok := result[fieldName].(bson.M)
			if !ok {
				http.Error(w, fieldName+" is not a JSON object", http.StatusInternalServerError)
				return
			}

			fieldValueAsBytes, err := json.Marshal(fieldValue)
			if err != nil {
				http.Error(w, fieldName+" could not be converted to string", http.StatusInternalServerError)
				return
			}

			fieldValueAsString := string(fieldValueAsBytes)
			db.SendOneField(w, fieldValueAsString, fieldName)
		},
	)
}

func GetJSONObjectByName(usersCollection *mongo.Collection, fieldName string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.D{
				{fieldName, 1},
			}
			findOptions := options.FindOne().SetProjection(projection)
			name := r.PathValue("name")
			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"name", name}}, findOptions)
			if err != nil {
				http.Error(w, "document with name "+name+" not found", http.StatusNotFound)
				return
			}

			fieldValue, ok := result[fieldName].(bson.M)

			if !ok {
				http.Error(w, fieldName+" is not a JSON object", http.StatusInternalServerError)
				return
			}

			fieldValueAsBytes, err := json.Marshal(fieldValue)
			if err != nil {
				http.Error(w, fieldName+" could not be converted to string", http.StatusInternalServerError)
				return
			}

			fieldValueAsString := string(fieldValueAsBytes)
			db.SendOneField(w, fieldValueAsString, fieldName)
		},
	)
}
