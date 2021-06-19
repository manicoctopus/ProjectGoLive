package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/authenticate", app.authenticateUser).Methods("POST")

	users := router.PathPrefix("/api/v1/users").Subrouter()
	users.HandleFunc("/create", app.createUser).Methods("POST")
	users.HandleFunc("/update", app.updateUser).Methods("PUT")
	users.HandleFunc("/delete", app.deleteUser).Methods("DELETE")
	users.HandleFunc("/{id}", app.retrieveUser).Methods("GET")
	users.HandleFunc("", app.retrieveAllUsers).Methods("GET")

	courses := router.PathPrefix("/api/v1/pdtsvcs").Subrouter()
	courses.HandleFunc("/create", app.createPdtsvc).Methods("POST")
	courses.HandleFunc("/update", app.updatePdtsvc).Methods("PUT")
	courses.HandleFunc("/delete", app.deletePdtsvc).Methods("DELETE")
	courses.HandleFunc("/{id}", app.retrievePdtsvc).Methods("GET")
	courses.HandleFunc("", app.retrieveAllPdtsvcs).Methods("GET")

	return router
}
