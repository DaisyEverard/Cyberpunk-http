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

	"main/config"
	"main/handlers"
)

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

	config.Collection = mongoClient.Database("cyberpunk-red").Collection("characters")

	defer mongoClient.Disconnect(context.TODO())

	http.HandleFunc("/document", handlers.WholeDocumentHandler)
	http.HandleFunc("/hp", handlers.HPHandler)

	fmt.Printf("\nServer is listening on port %s...", config.PortNumber)
	portAddress := ":" + config.PortNumber
	if err := http.ListenAndServe(portAddress, nil); err != nil {
		log.Fatal(err)
	}
}
