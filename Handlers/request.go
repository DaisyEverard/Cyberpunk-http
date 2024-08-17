package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}

	collection = mongoClient.Database("cyberpunk-red").Collection("characters")

	defer mongoClient.Disconnect(context.TODO())

	http.HandleFunc("/document", wholeDocumentHandler)
	http.HandleFunc("/hp", HPHandler)

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func wholeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWholeDocumentHandler(w, r)
	default:
		http.Error(w, "Method not allowed, GET only", http.StatusMethodNotAllowed)
	}
}

func HPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHPHandler(w, r)
	case http.MethodPost:
		updateHP(w, r)
	default:
		http.Error(w, "Method not allowed, GET and POST only", http.StatusMethodNotAllowed)
	}
}
