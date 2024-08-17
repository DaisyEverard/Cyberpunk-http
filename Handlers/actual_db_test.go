package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const serverPort = 1066 // Battle of Hastings
const testID = "66a0cbdd4f7c2fdd38865c26"
const testName = "Johnny Silverhand"

func connectToMongo() {
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

	fmt.Printf("\nServer is listening on port %v...", serverPort)
	address := fmt.Sprintf(":%d", serverPort)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}

// test is timing out before it can run
// if doing docker image, how to assert initial state?
func TestAll(t *testing.T) {
	connectToMongo()
	t.Run("getHPbyName", func(t *testing.T) {
		requestURL := fmt.Sprintf("http://localhost:%d/hp?id=%s", serverPort, testID)
		res, _ := http.Get(requestURL)

		resBody, _ := ioutil.ReadAll(res.Body)

		fmt.Printf("body: %v", resBody)

		// hp -  get body.hp
		// if !equal(hp, expectedHP) {
		// 	t.Errorf("got %v but wanted %v", got, want)
		// }
	})
}
