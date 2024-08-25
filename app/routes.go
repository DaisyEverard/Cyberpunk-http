package app

import (
	"net/http"

	"main/app/handlers"

	"go.mongodb.org/mongo-driver/mongo"
)

func addRoutes(mux *http.ServeMux, usersCollection *mongo.Collection) {
	// handle multiple docs with same name

	// DOCUMENT
	mux.Handle("GET /document/by_id/{id}", handlers.GetDocumentByID(usersCollection))
	mux.Handle("GET /document/by_name/{name}", handlers.GetDocumentByName(usersCollection))

	// INT64
	// HP
	mux.Handle("GET /hp/by_id/{id}", handlers.GetInt64ByID(usersCollection, "hp"))
	mux.Handle("GET /hp/by_name/{name}", handlers.GetInt64ByName(usersCollection, "hp"))
	mux.Handle("POST /hp/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "hp"))
	mux.Handle("POST /hp/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "hp"))
	// HUMANITY
	mux.Handle("GET /humanity/by_id/{id}", handlers.GetInt64ByID(usersCollection, "humanity"))
	mux.Handle("GET /humanity/by_name/{name}", handlers.GetInt64ByName(usersCollection, "humanity"))
	mux.Handle("POST /humanity/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "humanity"))
	mux.Handle("POST /humanity/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "humanity"))

	// STRING
	// NAME
	mux.Handle("GET /name/by_id/{id}", handlers.GetStringByID(usersCollection, "name"))
	mux.Handle("POST /name/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "name"))
	mux.Handle("POST /name/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "name"))
	// ROLE
	mux.Handle("GET /role/by_id/{id}", handlers.GetStringByID(usersCollection, "role"))
	mux.Handle("GET /role/by_name/{name}", handlers.GetStringByName(usersCollection, "role"))
	mux.Handle("POST /role/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "role"))
	mux.Handle("POST /role/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "role"))

	// OBJECTS
	// EFFECTS
	mux.Handle("GET /effects/by_id/{id}", handlers.GetJSONObjectByID(usersCollection, "effects"))
	mux.Handle("GET /effects/by_name/{name}", handlers.GetJSONObjectByName(usersCollection, "effects"))
	mux.Handle("POST /effects/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "effects"))
	mux.Handle("POST /effects/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "effects"))
	// SKILLS
	mux.Handle("GET /skills/by_id/{id}", handlers.GetJSONObjectByID(usersCollection, "skills"))
	mux.Handle("GET /skills/by_name/{name}", handlers.GetJSONObjectByName(usersCollection, "skills"))
	mux.Handle("POST /skills/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "skills"))
	mux.Handle("POST /skills/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "skills"))
	// STATS
	mux.Handle("GET /stats/by_id/{id}", handlers.GetJSONObjectByID(usersCollection, "stats"))
	mux.Handle("GET /stats/by_name/{name}", handlers.GetJSONObjectByName(usersCollection, "stats"))
	mux.Handle("POST /stats/by_id/{id}", handlers.UpdateFieldByID(usersCollection, "stats"))
	mux.Handle("POST /stats/by_name/{name}", handlers.UpdateFieldByName(usersCollection, "stats"))

	// OBJECTID
	// ID
	mux.Handle("GET /id/by_name/{name}", handlers.GetObjectIDByName(usersCollection))

	// DEFAULT
	mux.Handle("/{$}", http.NotFoundHandler())
}
