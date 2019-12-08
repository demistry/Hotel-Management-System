package main

import (
	"fmt"
	"github.com/Demistry/Hotel-Management-System/src/controllers/admin"
	"github.com/Demistry/Hotel-Management-System/src/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var GlobalRouter *mux.Router

func main(){
	fmt.Print("Starting server...\n")
	//mongoChannel := make
	err := godotenv.Load(utils.EnvironmentVariableFilename)
	if err != nil{
		fmt.Println("Error loading environment files ", err.Error())
	}
	//admin.SendMail()
	go admin.InitializeMongoDb()

	//for chann := range
	initializeRoutes()
}

func initializeRoutes(){
	GlobalRouter = mux.NewRouter()
	addRoutes()

	port, ok := os.LookupEnv("PORT")
	if ok == false {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":" + port, GlobalRouter))
}

