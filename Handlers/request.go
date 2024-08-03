package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"errors"
)

// CHARACTER STORES
type CharacterStore interface {
	FindOne(ctx context.Context, filter interface{}) (bson.M, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
}

type MongoCharacterStore struct {
	collection *mongo.Collection
}

func (m *MongoCharacterStore) FindOne(ctx context.Context, filter interface{}) (bson.M, error) {
	var result bson.M
	err := m.collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

func (m *MongoCharacterStore) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.collection.UpdateOne(ctx, filter, update)
}

type MockCharacterStore struct {
	Data []bson.M
}

func (m *MockCharacterStore) FindOne(ctx context.Context, filter interface{}) (bson.M, error) {
	filterMap, ok := filter.(bson.D)
	if !ok {
		return nil, errors.New("invalid filter type, expected bson.D")
	}

	for _, doc := range m.Data {
		matches := true
		for _, filterElem := range filterMap {
			intFilterValue, _ := strconv.Atoi(filterElem.Value.(string))
			if doc[filterElem.Key] != intFilterValue {
				matches = false
				break
			}
		}

		if matches {
			return doc, nil
		}
}

return nil, mongo.ErrNoDocuments
}

func (m *MockCharacterStore) UpdateOne(ctx context.Context, filter, update interface{}) (*mongo.UpdateResult, error) {
	filterMap, ok := filter.(bson.D)
	if !ok {
		return nil, errors.New("invalid filter type, expected bson.D")
	}

	fmt.Printf("update: %v", update)
	updateMap, ok := update.(bson.D)
	if !ok {
		return nil, errors.New("invalid update type, expected bson.D")
	}

	var matchedCount, modifiedCount int64

	for i, doc := range m.Data {
		matches := true
		for _, filterElem := range filterMap {
			if doc[filterElem.Key] != filterElem.Value {
				matches = false
				break
			}
		}

		if matches {
			matchedCount++
			updatedDoc := doc

			for _, updateElem := range updateMap {
				if updateElem.Key == "$set" {
					updates, ok := updateElem.Value.(bson.M)
					if !ok {
						return nil, errors.New("invalid update format for $set")
					}

					for k, v := range updates {
						updatedDoc[k] = v
					}

					m.Data[i] = updatedDoc
					modifiedCount++
				} else {
					return nil, errors.New("unsupported update operator")
				}
			}

			break
		}
	}

	if matchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &mongo.UpdateResult{
		MatchedCount:  matchedCount,
		ModifiedCount: modifiedCount,
	}, nil
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

	collection := mongoClient.Database("cyberpunk-red").Collection("characters")
	var store CharacterStore
	store = &MongoCharacterStore{collection}

	defer mongoClient.Disconnect(context.TODO());

	http.HandleFunc("/HP", makeHPHandler(store))

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func makeHPHandler(store CharacterStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHPHandler(w, r, store)
		case http.MethodPost:
			updateHPHandler(w, r, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getHPHandler(w http.ResponseWriter, r *http.Request, store CharacterStore) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	// result, err := store.FindOne(context.TODO(), bson.D{{"role", "Medtech"}}) // this worked
	// role is a string field
	// _id is of type ObjectID, user_id is an int32.
	// you need to convert to the right type to feed to the filter

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}
	result, err := store.FindOne(context.TODO(), bson.D{{"_id", objID}}) 

	if err == mongo.ErrNoDocuments {
		http.Error(w, fmt.Sprintf("No document found with the id %s", id), http.StatusNotFound)
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

func updateHPHandler(w http.ResponseWriter, r *http.Request, store CharacterStore) {
	var updateData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	id, ok := updateData["id"].(float64)
	if !ok {
		http.Error(w, "id field is required", http.StatusBadRequest)
		return
	}

	idStr := strconv.FormatFloat(id, 'f', 0, 64)

	update := bson.D{{"$set", updateData}}
	_, err := store.UpdateOne(context.TODO(), bson.D{{"id", idStr}}, update)
	if err == mongo.ErrNoDocuments {
		http.Error(w, fmt.Sprintf("No document found with the id %s", idStr), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HP of character %s updated successfully", idStr)
}