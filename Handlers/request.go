package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func getOneFromCollection(ctx context.Context, filter interface{}) (bson.M, error) {
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	// This is the bit it's failing on 
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}

	collection = mongoClient.Database("cyberpunk-red").Collection("characters")

	defer mongoClient.Disconnect(context.TODO());

	http.HandleFunc("/HP", HPHandler)

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func HPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDocumentHandler(w, r)
	case http.MethodPost:
		updateHP(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateHP(w http.ResponseWriter, r *http.Request) {
	var updateData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received updateData: %+v\n", updateData)

	id, _ := updateData["id"].(string)
	name, _ := updateData["name"].(string)

	if (id == "" && name == "") {
		http.Error(w, "name or id parameter is required", http.StatusBadRequest)
		return
	}

	newHP, ok := updateData["hp"] // int32 in Mongo
    if !ok {
        http.Error(w, "hp field is required and must be a number", http.StatusBadRequest)
        return
    }
	update := bson.D{{"$set", bson.M{"hp": newHP}}}

	var result *mongo.UpdateResult
	// Check if HP in document is already the same as the new one? 
	// result = {MatchedCount, ModifiedCount, UpsertedCount, UpsertedID}
	// MatchedCount is 0: No documents matched the filter.
	// ModifiedCount is 0: No documents were modified.
	// UpsertedCount is 0: No documents were upserted. (if no id exists, a new one is made with the value. Off by default)
	// UpsertedID is nil: No ID was generated because no upsert occurred.

	if id != "" {
		result = updateHPByID(w, r, update, id)
	} else if name != "" {
		result = updateHPByName(w,r,update, name)
	}

	if result.ModifiedCount == 0 {
        http.Error(w, "No document was updated. It may already have the desired value.", http.StatusNotModified)
        return
    }

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HP of character with id %s updated successfully", id)
}

func updateHPByID(w http.ResponseWriter, r *http.Request, update bson.D, id string) *mongo.UpdateResult {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return nil
	}

	result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", objID}}, update)

	if (err == mongo.ErrNoDocuments || result.MatchedCount == 0) {
		http.Error(w, fmt.Sprintf("No document found with the id %s", objID), http.StatusNotFound)
		return nil
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return result
}

func updateHPByName(w http.ResponseWriter, r *http.Request,update bson.D, name string) *mongo.UpdateResult {
	result, err := collection.UpdateOne(context.TODO(), bson.D{{"name", name}}, update)
	if (err == mongo.ErrNoDocuments || result.MatchedCount == 0) {
		http.Error(w, fmt.Sprintf("No document found with the name %s", name), http.StatusNotFound)
		return nil
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return result
}
