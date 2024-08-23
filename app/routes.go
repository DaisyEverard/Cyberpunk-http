package app

import (
	"net/http"

	"main/app/handlers"

	"go.mongodb.org/mongo-driver/mongo"
)

func addRoutes(mux *http.ServeMux, usersCollection *mongo.Collection) {

	mux.Handle("GET /document/id/{id}", handlers.GetDocumentByID(usersCollection))
	mux.Handle("GET /document/name/{name}", handlers.GetDocumentByName(usersCollection))

	mux.Handle("GET /hp/id/{id}", handlers.GetHPByID(usersCollection))
	mux.Handle("GET /hp/name/{name}", handlers.GetHPByName(usersCollection))
	// mux.Handler("GET /hp/id+name/{name, id}", handlers.GetUser(usersCollection))
	mux.Handle("POST /hp/id/{id}", handlers.UpdateHPByID(usersCollection))
	mux.Handle("POST /hp/name/{name}", handlers.UpdateHPByName(usersCollection))

	mux.Handle("/{$}", http.NotFoundHandler())
}
