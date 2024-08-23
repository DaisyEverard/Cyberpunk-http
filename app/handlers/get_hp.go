package handlers

import (
	"context"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/config"
	"main/app/db"
)

func GetHPHandler(w http.ResponseWriter, r *http.Request) {
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
		db.SendOneField(w, HPasString, "HP")
	} else if name != "" {
		HP := getHPByName(w, name)
		HPasString := strconv.FormatFloat(HP, 'f', -1, 64)
		db.SendOneField(w, HPasString, "HP")
	}
}

func getHPByID(w http.ResponseWriter, id string) float64 {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// This sends an HTTP response with the error message, no return value
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return config.UndefinedFloatValue
		// is further error handling necessary here, context to close?
	}

	projection := bson.D{
		{"hp", 1},
	}
	findOptions := options.FindOne().SetProjection(projection)

	result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"_id", objID}}, findOptions)
	if err != nil {
		http.Error(w, "HP not found", http.StatusNotFound)
		return config.UndefinedFloatValue
	}
	hp, ok := result["hp"].(float64)

	if !ok {
		http.Error(w, "Failed to convert HP to float64", http.StatusInternalServerError)
		return config.UndefinedFloatValue
	}

	return hp
}

func getHPByName(w http.ResponseWriter, name string) float64 {
	projection := bson.D{
		{"hp", 1},
	}
	findOptions := options.FindOne().SetProjection(projection)

	result, err := db.CallFindOneWithOptions(context.TODO(), bson.D{{"name", name}}, findOptions)
	if err != nil {
		http.Error(w, "HP not found", http.StatusNotFound)
		return config.UndefinedFloatValue
	}

	hp, ok := result["hp"].(float64)
	if !ok {
		http.Error(w, "Failed to convert HP to float64", http.StatusInternalServerError)
		return config.UndefinedFloatValue
	}

	return hp
}
