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

	pdtsvcs := router.PathPrefix("/api/v1/pdtsvcs").Subrouter()
	pdtsvcs.HandleFunc("/create", app.createPdtsvc).Methods("POST")
	pdtsvcs.HandleFunc("/update", app.updatePdtsvc).Methods("PUT")
	pdtsvcs.HandleFunc("/delete", app.deletePdtsvc).Methods("DELETE")
	pdtsvcs.HandleFunc("/{id}", app.retrievePdtsvc).Methods("GET")
	pdtsvcs.HandleFunc("", app.retrieveAllPdtsvcs).Methods("GET")

	listings := router.PathPrefix("/api/v1/listings").Subrouter()
	listings.HandleFunc("/create", app.createListing).Methods("POST")
	listings.HandleFunc("/update", app.updateListing).Methods("PUT")
	listings.HandleFunc("/delete", app.deleteListing).Methods("DELETE")
	listings.HandleFunc("/{id}", app.retrieveListing).Methods("GET")
	listings.HandleFunc("", app.retrieveAllListings).Methods("GET")
	listings.HandleFunc("/pdtsvcs/{id}", app.retrieveAllPdtsvcsByID).Methods("GET")
	listings.HandleFunc("/reviews/{id}", app.retrieveAllReviewsByID).Methods("GET")

	reviews := router.PathPrefix("/api/v1/reviews").Subrouter()
	reviews.HandleFunc("/create", app.createReview).Methods("POST")
	reviews.HandleFunc("/update", app.updateReview).Methods("PUT")
	reviews.HandleFunc("/delete", app.deleteReview).Methods("DELETE")
	reviews.HandleFunc("/{id}", app.retrieveReview).Methods("GET")
	reviews.HandleFunc("", app.retrieveAllReviews).Methods("GET")

	categories := router.PathPrefix("/api/v1/categories").Subrouter()
	categories.HandleFunc("/create", app.createCategory).Methods("POST")
	categories.HandleFunc("/update", app.updateCategory).Methods("PUT")
	categories.HandleFunc("/delete", app.deleteCategory).Methods("DELETE")
	categories.HandleFunc("/{id}", app.retrieveCategory).Methods("GET")
	categories.HandleFunc("", app.retrieveAllCategories).Methods("GET")

	return router
}
