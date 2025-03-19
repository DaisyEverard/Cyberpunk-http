package handlers

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/db"
)

func GetCharacterNames(usersCollection *mongo.Collection) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.M{"name": 1, "_id": 0}
			findOptions := options.Find().SetProjection(projection)

			cursor, err := db.CallFindWithOptions(context.TODO(), bson.M{}, findOptions)
			if err != nil {
				http.Error(w, "Error fetching data", http.StatusInternalServerError)
				return
			}
			defer cursor.Close(context.TODO())

			// Collect names in a slice
			var names []bson.M
			if err := cursor.All(context.TODO(), &names); err != nil {
				http.Error(w, "Error decoding data", http.StatusInternalServerError)
				return
			}

			db.SendNames(w, names)
		},
	)
}