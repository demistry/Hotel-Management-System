package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"encoding/json"
	"time"
)

//Collection parameters
const DatabaseName = "HotSysDatabase"
const HotelCollection = "hotels"
const EnvironmentVariableFilename = "src/app/EnvironmentVariables.env"
const HerokuBaseUrl = "https://hotsys.herokuapp.com/"
const ConfirmMailEndpoint = HerokuBaseUrl + "admin/create/confirm/"

const (
	MinHashCost = 4
	DefaultCost = 10
	MaxHashCost = 31
)

func GetHashedPassword(password string) string{
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte (password), DefaultCost)
	if err != nil{
		log.Print("Error in hashing ", err.Error())
		return ""
	}
	return string(hashedPasswordBytes)
}

func HandleError(statusCode int , genericResponse interface{}, err error,resp http.ResponseWriter){
	if err != nil{
		log.Println("Error that occurred is ", err.Error())
	}
	resp.WriteHeader(statusCode)
	json.NewEncoder(resp).Encode(genericResponse)
}

func GetHotelCollection(mongoClient *mongo.Client, uri string)(*mongo.Collection, context.Context, context.CancelFunc){
	collection := mongoClient.Database(DatabaseName).Collection(HotelCollection)
	mongoContext,cancel := context.WithTimeout(context.Background(), 28 * time.Second)
	return collection,mongoContext,cancel

}