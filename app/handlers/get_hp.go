package handlers

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/config"
	"main/app/db"
)

func GetHPByID(w http.ResponseWriter, id string) float64 {
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

func GetHPByName(w http.ResponseWriter, name string) float64 {
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
