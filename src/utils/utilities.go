package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetHotelCollection(mongoClient *mongo.Client, uri string)(*mongo.Collection, context.Context, context.CancelFunc){
	clientOptions := options.Client().ApplyURI(uri)
	mongoContext,cancel := context.WithTimeout(context.Background(), 28 * time.Second)
	mongoLocal,err := mongo.Connect(mongoContext, clientOptions)
	if err != nil{
		log.Println("Could not connect here....", err.Error())
		collection := mongoClient.Database(DatabaseName).Collection(HotelCollection)
		return collection, mongoContext, cancel
	}
	collection := mongoLocal.Database(DatabaseName).Collection(HotelCollection)

	return collection,mongoContext,cancel
}