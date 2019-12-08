package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

//Collection parameters
const DatabaseName = "HotelManagementSystemDatabase"
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

func GetHotelCollection(mongoClient *mongo.Client)(*mongo.Collection, context.Context, context.CancelFunc){
	collection := mongoClient.Database(DatabaseName).Collection(HotelCollection)
	mongoContext,cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	return collection,mongoContext,cancel
}