package main

import (
	"fmt"
	"github.com/Demistry/Hotel-Management-System/Source_Files/controllers/admin"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var GlobalRouter *mux.Router

func main(){
	fmt.Print("Starting server...\n")
	admin.InitializeMongoDb()
	initializeRoutes()
}

func initializeRoutes(){
	GlobalRouter = mux.NewRouter()
	addRoutes()
	log.Fatal(http.ListenAndServe(":12345", GlobalRouter))
}

