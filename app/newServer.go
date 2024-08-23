package app

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewServer(usersCollection *mongo.Collection) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, usersCollection)

	var handler http.Handler = mux

	return handler
}