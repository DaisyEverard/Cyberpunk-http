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

func GetObjectIDByName(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.D{
				{"_id", 1},
			}
			findOptions := options.FindOne().SetProjection(projection)

			name := r.PathValue("name")
			result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"name", name}}, findOptions)
			if err != nil {
				http.Error(w, "id not found", http.StatusNotFound)
				return
			}

			fieldValue, ok := result["_id"].(primitive.ObjectID)
			if !ok {
				http.Error(w, "Failed to convert id to ObjectID", http.StatusInternalServerError)
				return
			}

			fieldValueStr := fieldValue.Hex()
			db.SendOneField(w, fieldValueStr, "id")
		},
	)
}
