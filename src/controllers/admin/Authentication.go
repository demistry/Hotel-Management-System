package admin

import (
	"context"
	"encoding/json"
	"github.com/Demistry/Hotel-Management-System/src/models"
	"github.com/Demistry/Hotel-Management-System/src/responses"
	"github.com/Demistry/Hotel-Management-System/src/utils"
	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)


var mongoClient *mongo.Client


func InitializeMongoDb(){
	mongoContext,_ := context.WithTimeout(context.Background(), 15 * time.Second)
	uri,ok := os.LookupEnv("MONGODB_URI")
	if ok == false{
		uri = "mongodb://localhost:27017"
	}
	clientOptions := options.Client().ApplyURI(uri)
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
	collection, mongoContext, cancel := utils.GetHotelCollection(mongoClient)
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
	adminUser.IsUserVerified = false
	adminUser.CreatedAt = time.Now()
	adminUser.LinkExpiresAt = time.Now().Add(30 * time.Second)
	insertedID,er := collection.InsertOne(mongoContext, &adminUser)
	if er != nil{
		response.WriteHeader(http.StatusInternalServerError)
		errResponse := responses.GenericResponse{Status:false, Message:"Internal Server Error"}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	//sendMail(adminUser.HotelEmail, adminUser.HotelName, insertedID.InsertedID.(primitive.ObjectID).Hex())
	json.NewEncoder(response).Encode(insertedID)
}

func VerifyAdminEmail(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	idParameter := mux.Vars(request)
	id,_ := primitive.ObjectIDFromHex(idParameter["id"])
	var admin models.AdminUser
	collection, mongoContext, cancel := utils.GetHotelCollection(mongoClient)
	defer cancel()
	filter := bson.M{"_id": id}
	updateFilter := bson.M{"$set": bson.M{"isUserVerified": true}}
	err := collection.FindOne(mongoContext, filter).Decode(&admin)
	if err != nil{
		response.WriteHeader(http.StatusOK)
		errResponse := responses.GenericResponse{Status:false, Message:"Could not find user"}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	if admin.LinkExpiresAt.After(time.Now()){
		if !admin.IsUserVerified{
			_, _ = collection.UpdateOne(mongoContext, filter, updateFilter)
			response.WriteHeader(http.StatusOK)
			errResponse := responses.GenericResponse{Status:true, Message:"User email successfully verified"}
			json.NewEncoder(response).Encode(errResponse)
			return
		} else{
			response.WriteHeader(http.StatusOK)
			errResponse := responses.GenericResponse{Status:false, Message:"User is already verified"}
			json.NewEncoder(response).Encode(errResponse)
			return
		}
	}else{
		response.WriteHeader(http.StatusOK)
		errResponse := responses.GenericResponse{Status:false, Message:"Verification link expired"}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
}


func sendMail(emailAddress string, username string, userId string){
	from := mail.NewEmail("HotSys", "Hotsys@mail.com")
	subject := "Email Verification for HotSys"
	to := mail.NewEmail(username, emailAddress)
	content := mail.NewContent("text/plain", "Click on the link below to verify your email address for " + username + "\n " + utils.HerokuBaseUrl + utils.ConfirmMailEndpoint + userId)
	m := mail.NewV3MailInit(from, subject, to, content)
	apiKey,ok := os.LookupEnv("SENDGRID_API_KEY")
	if ok == false{
		apiKey = os.Getenv("SENDGRID_API_KEY")
	}
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	}
}




