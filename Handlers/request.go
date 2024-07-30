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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. See: www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("cyberpunk-red").Collection("characters")
}

func main() {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	http.HandleFunc("/HP", HPHandler)

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// get
func HPHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getHPHandler(w,r)
	case http.MethodPost:
		updateHPHandler(w,r)
	default: 
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHPHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		http.Error(w, fmt.Sprintf("No document found with the name %s", name), http.StatusNotFound)
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

func updateHPHandler(w http.ResponseWriter, r *http.Request) {
	var updateData bson.M
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	name, ok := updateData["name"].(string)
	if !ok || name == "" {
		http.Error(w, "Name field is required", http.StatusBadRequest)
		return
	}

	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), bson.D{{"name", name}}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HP of character %s updated successfully", name)
}