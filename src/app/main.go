package main

import (
	"fmt"
	"github.com/Demistry/Hotel-Management-System/src/controllers/admin"
	"github.com/Demistry/Hotel-Management-System/src/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var GlobalRouter *mux.Router

func main(){
	fmt.Print("Starting server...\n")
	err := godotenv.Load(utils.EnvironmentVariableFilename)
	if err != nil{
		fmt.Println("Error loading environment files ", err.Error())
	}
	//admin.SendMail()
	admin.InitializeMongoDb()
	initializeRoutes()
}

func initializeRoutes(){
	GlobalRouter = mux.NewRouter()
	addRoutes()
	log.Fatal(http.ListenAndServe(":12345", GlobalRouter))
}

