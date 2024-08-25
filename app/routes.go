package app

import (
	"net/http"

	"main/app/handlers"

	"go.mongodb.org/mongo-driver/mongo"
)

func addRoutes(mux *http.ServeMux, usersCollection *mongo.Collection) {

	mux.Handle("GET /document/id/{id}", handlers.GetDocumentByID(usersCollection))
	mux.Handle("GET /document/name/{name}", handlers.GetDocumentByName(usersCollection))

	mux.Handle("GET /hp/id/{id}", handlers.GetNumberByID(usersCollection, "hp"))
	mux.Handle("GET /hp/name/{name}", handlers.GetNumberByName(usersCollection, "hp"))
	// mux.Handler("GET /hp/id+name/{name, id}", handlers.GetUser(usersCollection))
	mux.Handle("POST /hp/id/{id}", handlers.UpdateFieldByID(usersCollection, "hp"))
	mux.Handle("POST /hp/name/{name}", handlers.UpdateFieldByName(usersCollection, "hp"))

	mux.Handle("/{$}", http.NotFoundHandler())
}

// effects: json objects
// humanity: int32
// name: string
// role: string
// skills: json object
// stats: object
// userID: ObjectID
