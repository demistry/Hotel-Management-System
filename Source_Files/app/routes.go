package main

import (
	"github.com/Demistry/Hotel-Management-System/Source_Files/controllers/admin"
	"github.com/gorilla/mux"
)

var adminRoutes *mux.Router
func addRoutes(){
	adminRoutes = GlobalRouter.PathPrefix("/admin").Subrouter()
	adminRoutes.HandleFunc("/create",admin.CreateNewHotelAdmin).Methods("POST")
}
