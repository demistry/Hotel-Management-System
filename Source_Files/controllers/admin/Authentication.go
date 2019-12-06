package admin

import (
	"context"
	"encoding/json"
	"github.com/Demistry/Hotel-Management-System/Source_Files/models"
	"github.com/Demistry/Hotel-Management-System/Source_Files/responses"
	"github.com/Demistry/Hotel-Management-System/Source_Files/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"regexp"
	"time"
)


var mongoClient *mongo.Client


func InitializeMongoDb(){
	mongoContext,_ := context.WithTimeout(context.Background(), 15 * time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient,_ = mongo.Connect(mongoContext, clientOptions)
}

func CreateNewHotelAdmin(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	var adminUser *models.AdminUser
	err := json.NewDecoder(request.Body).Decode(&adminUser)
	if err != nil{
		response.WriteHeader(http.StatusForbidden)
		errResponse := responses.GenericResponse{Status:false, Message:"Missing field(s)"}
		log.Print("Error in decoding body is ", err.Error())
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	collection := mongoClient.Database(utils.DatabaseName).Collection(utils.HotelCollection)
	mongoContext,cancel := context.WithTimeout(context.Background(), 8 * time.Second)
	defer cancel()
	if isEmailValid,_ := regexp.MatchString("(\\w+)@(\\w+)\\.com", adminUser.HotelEmail);!isEmailValid{
		response.WriteHeader(http.StatusOK)
		errResponse := responses.GenericResponse{Status: false, Message: "Email:" + adminUser.HotelEmail + " is not a valid email.."}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	adminUser.HotelPassword = utils.GetHashedPassword(adminUser.HotelPassword)
	filter := bson.M{"hotelEmail": adminUser.HotelEmail}
	findError := collection.FindOne(mongoContext, filter).Decode(&adminUser)

	if findError == nil { //check if database already contains email
		response.WriteHeader(http.StatusOK)
		errResponse := responses.GenericResponse{Status: false, Message: "Email:" + adminUser.HotelEmail + " already in use."}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	insertedAdmin,er := collection.InsertOne(mongoContext, &adminUser)
	if er != nil{
		response.WriteHeader(http.StatusInternalServerError)
		errResponse := responses.GenericResponse{Status:false, Message:"Internal Server Error"}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	json.NewEncoder(response).Encode(insertedAdmin)
}




