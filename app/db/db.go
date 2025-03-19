package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/config"
)

// CALL
func CallFindWithOptions(ctx context.Context, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := config.Collection.Find(context.TODO(), filter, options)
	return cursor, err
}

func CallFindOne(ctx context.Context, filter interface{}) (bson.M, error) {
	var result bson.M
	err := config.Collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func CallFindOneWithOptions(ctx context.Context, filter interface{}, options *options.FindOneOptions) (bson.M, error) {
	var result bson.M
	err := config.Collection.FindOne(context.TODO(), filter, options).Decode(&result)
	return result, err
}

// SEND
func SendDocument(w http.ResponseWriter, result bson.M, err error, queryType string) {
	if err == mongo.ErrNoDocuments {
		http.Error(w, fmt.Sprintf("No document found with that %s", queryType), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SendOneField(w http.ResponseWriter, value string, queryType string) {
	result := "{fieldName:" + queryType + ",value:" + value + "}"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
