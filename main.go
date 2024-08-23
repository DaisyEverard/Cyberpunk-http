package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"main/app"
	"main/app/config"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	MONGODB_USERNAME := os.Getenv("MONGODB_USERNAME")
	MONGODB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	PORT := os.Getenv("PORT")

	connection_string := strings.Replace(MONGODB_URI, "<db_username>", MONGODB_USERNAME, 1)
	connection_string = strings.Replace(connection_string, "<db_password>", MONGODB_PASSWORD, 1)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connection_string))
	if err != nil {
		log.Fatal(err)
	}

	config.Collection = mongoClient.Database("cyberpunk-red").Collection("characters")

	defer mongoClient.Disconnect(context.TODO())

	srv := app.NewServer(config.Collection)
	httpServer := &http.Server{
		Addr:    `:` + PORT,
		Handler: srv,
	}

	fmt.Printf("\nServer is listening on port %s...", PORT)
	httpServer.ListenAndServe()
}
