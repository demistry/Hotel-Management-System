package main

import (
	"github.com/Demistry/Hotel-Management-System/src/controllers/admin"
	"github.com/Demistry/Hotel-Management-System/src/controllers/users"
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
	adminRoutes.HandleFunc("/create/confirm/{id}", admin.VerifyAdminEmail).Methods("GET")
	adminRoutes.HandleFunc("/login", admin.LoginUser).Methods("POST")
	adminRoutes.HandleFunc("/resetPassword", admin.ResendVerificationMailForResetPassword).Methods("POST")
	adminRoutes.HandleFunc("/resetUserPassword", admin.ResetPassword).Methods("POST")


	userRoutes := GlobalRouter.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("/create", users.CreateUser).Methods("POST")
}
