package handlers

import (
	"context"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/db"
)

func GetHPByID(usersCollection *mongo.Collection) http.HandlerFunc {
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
				{"hp", 1},
			}
			findOptions := options.FindOne().SetProjection(projection)

			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"_id", objID}}, findOptions)
			if err != nil {
				http.Error(w, "HP not found", http.StatusNotFound)
				return
			}
			hp, ok := result["hp"].(float64)

			if !ok {
				http.Error(w, "Failed to convert HP to float64", http.StatusInternalServerError)
				return
			}

			HPasString := strconv.FormatFloat(hp, 'f', -1, 64)
			db.SendOneField(w, HPasString, "HP")
		},
	)
}

func GetHPByName(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.D{
				{"hp", 1},
			}
			findOptions := options.FindOne().SetProjection(projection)

			name := r.PathValue("name")
			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"name", name}}, findOptions)
			if err != nil {
				http.Error(w, "HP not found", http.StatusNotFound)
				return
			}

			hp, ok := result["hp"].(float64)
			if !ok {
				http.Error(w, "Failed to convert HP to float64", http.StatusInternalServerError)
				return
			}

			HPasString := strconv.FormatFloat(hp, 'f', -1, 64)
			db.SendOneField(w, HPasString, "HP")
		},
	)
}
