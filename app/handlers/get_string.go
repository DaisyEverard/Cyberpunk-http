package handlers

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/db"
)

func GetStringByID(usersCollection *mongo.Collection, fieldName string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				// This sends an HTTP response with the error message, no return value
				http.Error(w, "Invalid id format", http.StatusBadRequest)
				return
				// is further error handling necessary here, context to close?
			}

			projection := bson.D{
				{fieldName, 1},
			}
			findOptions := options.FindOne().SetProjection(projection)

			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"_id", objID}}, findOptions)
			if err != nil {
				http.Error(w, "ID not found", http.StatusNotFound)
				return
			}
			fieldValue, ok := result[fieldName].(string)

			if !ok {
				http.Error(w, "Failed to convert "+fieldName+" to string", http.StatusInternalServerError)
				return
			}

			db.SendOneField(w, fieldValue, fieldName)
		},
	)
}

func GetStringByName(usersCollection *mongo.Collection, fieldName string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.D{
				{fieldName, 1},
			}
			findOptions := options.FindOne().SetProjection(projection)

			name := r.PathValue("name")
			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"name", name}}, findOptions)
			if err != nil {
				http.Error(w, fieldName+" not found", http.StatusNotFound)
				return
			}

			fieldValue, ok := result[fieldName].(string)
			if !ok {
				http.Error(w, "Failed to convert "+fieldName+" to string", http.StatusInternalServerError)
				return
			}

			db.SendOneField(w, fieldValue, fieldName)
		},
	)
}
