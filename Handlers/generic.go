package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const undefinedFloatValue = -999

func callFindOne(ctx context.Context, filter interface{}) (bson.M, error) {
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func callFindOneWithOptions(ctx context.Context, filter interface{}, options *options.FindOneOptions) (bson.M, error) {
	var result bson.M
	err := collection.FindOne(context.TODO(), filter, options).Decode(&result)
	return result, err
}

func sendDocument(w http.ResponseWriter, result bson.M, err error, queryType string) {
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

func sendStringField(w http.ResponseWriter, value string, queryType string) {
	result := "{fieldName:" + queryType + ",value:" + value + "}"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func sendFloat64Field(w http.ResponseWriter, err error, value float64, queryType string) {
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	result := "{fieldName:" + queryType + ",value:" + value + "}"

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(result); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
