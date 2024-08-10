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
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	http.HandleFunc("/HP", makeDocumentHandler(store))

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func makeDocumentHandler(store CharacterStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDocumentHandler(w, r, store)
		case http.MethodPost:
			updateHP(w, r, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func updateHP(w http.ResponseWriter, r *http.Request, store CharacterStore) {
	var updateData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received updateData: %+v\n", updateData)

	id, _ := updateData["id"].(string)
	name, _ := updateData["name"].(string)
	print("name: %v", name)
	print("id: %v", id)

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
		result = updateHPByID(w, r, store, update, id)
	} else if name != "" {
		result = updateHPByName(w,r, store,update, name)
	}

	if result.ModifiedCount == 0 {
        http.Error(w, "No document was updated. It may already have the desired value.", http.StatusNotModified)
        return
    }

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HP of character with id %s updated successfully", id)
}
 
func updateHPByID(w http.ResponseWriter, r *http.Request, store CharacterStore, update bson.D, id string) *mongo.UpdateResult {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return nil
	}

	result, err := store.UpdateOne(context.TODO(), bson.D{{"_id", objID}}, update)
	if (err == mongo.ErrNoDocuments || result.MatchedCount == 0) {
		http.Error(w, fmt.Sprintf("No document found with the id %s", objID), http.StatusNotFound)
		return nil
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return result
}

func updateHPByName(w http.ResponseWriter, r *http.Request, store CharacterStore,update bson.D, name string) *mongo.UpdateResult {
	fmt.Printf("\nname: %v\n", name)

	result, err := store.UpdateOne(context.TODO(), bson.D{{"name", name}}, update)
	if (err == mongo.ErrNoDocuments || result.MatchedCount == 0) {
		http.Error(w, fmt.Sprintf("No document found with the name %s", name), http.StatusNotFound)
		return nil
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return result
}