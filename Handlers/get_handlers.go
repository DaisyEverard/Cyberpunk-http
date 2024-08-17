package main

import (
	"context"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getWholeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	if name == "" && id == "" {
		http.Error(w, "name or id parameter is required", http.StatusBadRequest)
		return
	}

	// logic for only one query field at a time
	// implement multiple query fields at a time
	if id != "" {
		result, err := getDocumentByID(w, id)
		sendDocument(w, result, err, "id")
	} else if name != "" {
		result, err := getDocumentByName(w, name)
		sendDocument(w, result, err, "name")
	}
}

func getHPHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	if name == "" && id == "" {
		http.Error(w, "name or id parameter is required", http.StatusBadRequest)
		return
	}

	// logic for only one query field at a time
	// implement multiple query fields at a time
	if id != "" {
		HP := getHPByID(w, id)
		HPasString := strconv.FormatFloat(HP, 'f', -1, 64)
		sendStringField(w, HPasString, "HP")
	} else if name != "" {
		HP := getHPByName(w, name)
		HPasString := strconv.FormatFloat(HP, 'f', -1, 64)
		sendStringField(w, HPasString, "HP")
	}
}

func getDocumentByID(w http.ResponseWriter, id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// This sends an HTTP response with the error message, no return value
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return nil, err
		// is further error handling necessary here, context to close?
	}
	return callFindOne(context.TODO(), bson.D{{"_id", objID}})
}

func getDocumentByName(w http.ResponseWriter, name string) (bson.M, error) {
	return callFindOne(context.TODO(), bson.D{{"name", name}})
}

func getHPByID(w http.ResponseWriter, id string) float64 {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// This sends an HTTP response with the error message, no return value
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return undefinedFloatValue
		// is further error handling necessary here, context to close?
	}

	projection := bson.D{
		{"hp", 1},
	}
	findOptions := options.FindOne().SetProjection(projection)

	result, err := callFindOneWithOptions(context.TODO(), bson.D{{"_id", objID}}, findOptions)
	if err != nil {
		http.Error(w, "HP not found", http.StatusNotFound)
		return undefinedFloatValue
	}
	hp, ok := result["hp"].(float64)

	if !ok {
		http.Error(w, "Failed to convert HP to float64", http.StatusInternalServerError)
		return undefinedFloatValue
	}

	return hp
}

func getHPByName(w http.ResponseWriter, name string) float64 {
	projection := bson.D{
		{"hp", 1},
	}
	findOptions := options.FindOne().SetProjection(projection)

	result, err := callFindOneWithOptions(context.TODO(), bson.D{{"name", name}}, findOptions)
	if err != nil {
		http.Error(w, "HP not found", http.StatusNotFound)
		return undefinedFloatValue
	}

	hp, ok := result["hp"].(float64)
	if !ok {
		http.Error(w, "Failed to convert HP to float64", http.StatusInternalServerError)
		return undefinedFloatValue
	}

	return hp
}
