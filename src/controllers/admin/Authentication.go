package admin

import (
	"context"
	"encoding/json"
	"github.com/Demistry/Hotel-Management-System/src/models"
	"github.com/Demistry/Hotel-Management-System/src/responses"
	"github.com/Demistry/Hotel-Management-System/src/utils"
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
	insertedID,er := collection.InsertOne(mongoContext, &adminUser)
	if er != nil{
		response.WriteHeader(http.StatusInternalServerError)
		errResponse := responses.GenericResponse{Status:false, Message:"Internal Server Error"}
		json.NewEncoder(response).Encode(errResponse)
		return
	}
	sendMail(adminUser.HotelEmail, adminUser.HotelName, insertedID.InsertedID.(primitive.ObjectID).Hex())
	json.NewEncoder(response).Encode(insertedID)
}


func sendMail(emailAddress string, username string, userid string){
	from := mail.NewEmail("HotSys", "Hotsys@mail.com")
	subject := "Email Verification for HotSys"
	to := mail.NewEmail(username, emailAddress)
	content := mail.NewContent("text/plain", "Click on the link below to verify your email address\n " + utils.HerokuBaseUrl + utils.ConfirmMailEndpoint + userid)
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




