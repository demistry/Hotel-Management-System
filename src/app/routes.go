package main

import (
	"github.com/Demistry/Hotel-Management-System/src/controllers/admin"
	"github.com/gorilla/mux"
	"net/http"
)

var adminRoutes *mux.Router
func addRoutes(){
	GlobalRouter.HandleFunc("/", func(resp http.ResponseWriter, r *http.Request){
		_, _ = resp.Write([]byte("WELCOME TO HOTSYS APP"))
	})
	adminRoutes = GlobalRouter.PathPrefix("/admin").Subrouter()
	adminRoutes.HandleFunc("/create",admin.CreateNewHotelAdmin).Methods("POST")
}
