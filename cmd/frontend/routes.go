package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {

	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/", app.home).Methods("GET")

	employee := router.PathPrefix("/pdtsvc").Subrouter()
	employee.HandleFunc("/addForm", app.addPdtsvcForm).Methods("GET")
	employee.HandleFunc("/editForm/{id}", app.editPdtsvcForm).Methods("GET")
	employee.HandleFunc("/create", app.addPdtsvc).Methods("POST")
	employee.HandleFunc("/update", app.updatePdtsvc).Methods("POST")
	employee.HandleFunc("/delete/{id}", app.deletePdtsvc).Methods("GET")
	employee.HandleFunc("/{id}", app.viewPdtsvc).Methods("GET")

	router.HandleFunc("/pdtsvcs", app.viewAllPdtsvcs).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", app.fileServer))

	return router
}
