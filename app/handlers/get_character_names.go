package handlers

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/db"
)

func GetCharacterNamesAndIDs(usersCollection *mongo.Collection) http.HandlerFunc {
	// returns data in format [{"name": "name1"},{"name": "name2"},...]
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			projection := bson.M{ "_id": 1, "name": 1}
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