package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/app/config"
	"main/app/handlers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	MONGODB_USERNAME := os.Getenv("MONGODB_USER")
	MONGODB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	connection_string := strings.Replace(MONGODB_URI,"<db_username>",MONGODB_USERNAME, 1)
	connection_string = strings.Replace(connection_string,"<db_password>",MONGODB_PASSWORD, 1)

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connection_string))
	if err != nil {
		log.Fatal(err)
	}

	config.Collection = mongoClient.Database("cyberpunk-red").Collection("characters")

	defer mongoClient.Disconnect(context.TODO())

	http.HandleFunc("/document", handlers.WholeDocumentHandler)
	http.HandleFunc("/hp", handlers.HPHandler,)

	fmt.Printf("\nServer is listening on port %s...", config.PortNumber)
	portAddress := ":" + config.PortNumber
	if err := http.ListenAndServe(portAddress, nil); err != nil {
		log.Fatal(err)
	}
}