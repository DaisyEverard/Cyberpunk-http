package app

import (
	"net/http"

	"main/app/handlers"

	"go.mongodb.org/mongo-driver/mongo"
)

func addRoutes(mux *http.ServeMux, usersCollection *mongo.Collection) {
	// handle multiple docs with same name

	// DOCUMENT
	mux.Handle("GET /document/id/{id}", handlers.GetDocumentByID(usersCollection))
	mux.Handle("GET /document/name/{name}", handlers.GetDocumentByName(usersCollection))

	// NUMBERS
	// HP
	mux.Handle("GET /hp/id/{id}", handlers.GetInt64ByID(usersCollection, "hp"))
	mux.Handle("GET /hp/name/{name}", handlers.GetInt64ByName(usersCollection, "hp"))
	mux.Handle("POST /hp/id/{id}", handlers.UpdateFieldByID(usersCollection, "hp"))
	mux.Handle("POST /hp/name/{name}", handlers.UpdateFieldByName(usersCollection, "hp"))
	// HUMANITY
	mux.Handle("GET /humanity/id/{id}", handlers.GetInt64ByID(usersCollection, "humanity"))
	mux.Handle("GET /humanity/name/{name}", handlers.GetInt64ByName(usersCollection, "humanity"))
	mux.Handle("POST /humanity/id/{id}", handlers.UpdateFieldByID(usersCollection, "humanity"))
	mux.Handle("POST /humanity/name/{name}", handlers.UpdateFieldByName(usersCollection, "humanity"))

	// effects: json objects
	// name: string
	// role: string
	// skills: json object
	// stats: object
	// userID: ObjectID
	// DEFAULT
	mux.Handle("/{$}", http.NotFoundHandler())
}
