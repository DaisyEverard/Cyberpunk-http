package app

import (
	"net/http"

	"github.com/george-hanks/go-mongo-http-server/app/handlers"

	"go.mongodb.org/mongo-driver/mongo"
)

func addRoutes(mux *http.ServeMux, usersCollection *mongo.Collection) {

	mux.Handle("GET /document/id/{id}", handlers.GetHPbyID(usersCollection))
	mux.Handle("GET /document/name/{name}", handlers.GetHPbyName(usersCollection))


	mux.Handle("GET /hp/id/{id}", handlers.GetUser(usersCollection))
	mux.Handler("GET /hp/name/{name}", handlers.GetUser(usersCollection))
	// mux.Handler("GET /hp/id+name/{name, id}", handlers.GetUser(usersCollection))
	mux.Handle("POST /hp/id/{id}", handlers.GetUser(usersCollection))
	mux.Handler("POST /hp/name/{name}", handlers.GetUser(usersCollection))

	mux.Handle("/", http.NotFoundHandler())
}
