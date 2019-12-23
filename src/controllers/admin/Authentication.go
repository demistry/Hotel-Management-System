package admin

import (
	"context"
	"encoding/json"
	"github.com/Demistry/Hotel-Management-System/src/models"
	"github.com/Demistry/Hotel-Management-System/src/responses"
	"github.com/Demistry/Hotel-Management-System/src/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)


var mongoClient *mongo.Client
var uri string


func InitializeMongoDb(){


	mongoContext,_ := context.WithTimeout(context.Background(), 15 * time.Second)
	uriHere,ok := os.LookupEnv("MLAB_URI")
	uri = uriHere
	if ok == false{
		log.Println("Did not see uri from environment")
		uri = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(uri)
	mongoLocal,err := mongo.Connect(mongoContext, clientOptions)
	if err != nil{
		log.Println("Error with connecting to BD is ", err.Error())
		return
	}
	log.Println("Mongo db connected")
	mongoClient = mongoLocal
}

func CreateNewHotelAdmin(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	var adminUser *models.AdminUser
	err := json.NewDecoder(request.Body).Decode(&adminUser)
	if err != nil{
		log.Print("Error in decoding body is ", err.Error())
		utils.HandleError(http.StatusInternalServerError, responses.GenericResponse{Status:false, Message:"Missing field(s)"},err, response)
		return
	}
	collection, mongoContext, cancel := utils.GetHotelCollection(mongoClient,uri)
	defer cancel()
	defer mongoContext.Done()
	if isEmailValid,_ := regexp.MatchString("(\\w+)@(\\w+)\\.com", adminUser.HotelEmail);!isEmailValid {
		utils.HandleError(http.StatusOK, responses.GenericResponse{Status: false, Message: "Email:" + adminUser.HotelEmail + " is not a valid email.."},nil, response)
		return
	}
	adminUser.HotelPassword = utils.GetHashedPassword(adminUser.HotelPassword)
	filter := bson.M{"hotelEmail": adminUser.HotelEmail}
	findErrorChan := make(chan error)
	defer close(findErrorChan)
	go func(){
		findErrorChan <- collection.FindOne(mongoContext, filter).Decode(&adminUser)
	}()
	findError := <- findErrorChan
	if findError == nil { //check if database already contains email
		utils.HandleError(http.StatusOK, responses.GenericResponse{Status: false, Message: "Email:" + adminUser.HotelEmail + " already in use."},nil, response)
		return
	}
	adminUser.IsUserVerified = false
	adminUser.CreatedAt = time.Now()
	adminUser.LinkExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	insertionChannel := make(chan models.InsertionStruct)
	defer close(insertionChannel)
	go func() {
		insertionId, er := collection.InsertOne(mongoContext, &adminUser)
		insertionChannel <- models.InsertionStruct{
			InsertedId: insertionId,
			Er:         er,
		}
	}()
	insertedStruct := <- insertionChannel
	if insertedStruct.Er != nil{
		utils.HandleError(http.StatusInternalServerError, responses.GenericResponse{Status:false, Message:"Internal Server Error"},insertedStruct.Er, response)
		return
	}

	go sendMail(adminUser.HotelEmail, adminUser.HotelName, insertedStruct.InsertedId.InsertedID.(primitive.ObjectID).Hex())

	json.NewEncoder(response).Encode(responses.SuccessfulResponse{
		Status:  true,
		Message: "Successfully created account",
		Data:    insertedStruct.InsertedId.InsertedID,
	})
}

func VerifyAdminEmail(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	idParameter := mux.Vars(request)
	id,_ := primitive.ObjectIDFromHex(idParameter["id"])
	var admin models.AdminUser
	collection, mongoContext, cancel := utils.GetHotelCollection(mongoClient, uri)
	defer cancel()
	filter := bson.M{"_id": id}
	updateFilter := bson.M{"$set": bson.M{"isUserVerified": true}}
	channel := make(chan error)
	go func() {
		channel <- collection.FindOne(mongoContext, filter).Decode(&admin)
	}()
	err := <- channel
	if err != nil{
		utils.HandleError(http.StatusOK, responses.GenericResponse{Status:false, Message:"Could not find users"},err, response)
		return
	}
	if admin.LinkExpiresAt.After(time.Now()){
		if !admin.IsUserVerified{
			_, _ = collection.UpdateOne(mongoContext, filter, updateFilter)
			utils.HandleError(http.StatusOK, responses.GenericResponse{Status:true, Message:"User email successfully verified"},nil, response)
			return
		} else{
			utils.HandleError(http.StatusOK, responses.GenericResponse{Status:false, Message:"User is already verified"},nil, response)
			return
		}
	}else{
		utils.HandleError(http.StatusOK, responses.GenericResponse{Status:false, Message:"Verification link expired"},nil, response)
		return
	}
}

func LoginUser(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type","application/json")
	var loginObject *models.LoginRequest
	var adminUser *models.AdminUser
	err := json.NewDecoder(request.Body).Decode(&loginObject)
	if err != nil{
		utils.HandleError(http.StatusForbidden, responses.GenericResponse{Status:false, Message:"Missing field(s)"},err, response)
		return
	}

	filter := bson.M{"hotelEmail":loginObject.Email}
	collection, ctx, ctxCancel := utils.GetHotelCollection(mongoClient, uri)
	channel := make(chan error)
	defer ctxCancel()
	go func() {
		channel <- collection.FindOne(ctx,filter).Decode(&adminUser)
	}()
	for errorss := range channel{
		log.Println("result received from go routine")
		if errorss != nil{
			utils.HandleError(http.StatusOK, responses.GenericResponse{Status:false, Message:"Could not find users"},errorss, response)
			return
		}
		if !adminUser.IsUserVerified{
			utils.HandleError(http.StatusOK, responses.GenericResponse{Status:false, Message:"User is unverified"},errorss, response)
			return
		}

		isMatchedError := bcrypt.CompareHashAndPassword([]byte(adminUser.HotelPassword), []byte(loginObject.Password))
		if isMatchedError == nil{
			utils.HandleError(http.StatusOK, responses.SuccessfulResponse{Status:true, Message:"Successfully logged in", Data:adminUser.CreateResponse()},isMatchedError, response)
			return
		} else{
			utils.HandleError(http.StatusOK, responses.GenericResponse{Status:false, Message:"User password is incorrect."},isMatchedError, response)
			return
		}
	}

}

func ForgottenPassword(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	err := request.ParseForm()
	if err != nil{
		utils.HandleError(http.StatusForbidden,responses.GenericResponse{Status:false, Message:"Missing field(s)"},err, response)
		return
	}
	keyValues := request.Form
	email := keyValues.Get("email")
	log.Println("The email received for forgotten password is ", email)
	json.NewEncoder(response).Encode(responses.GenericResponse{
		Status:  false,
		Message: "Email recieved is " + email,
	})
}


func sendMail(emailAddress string, username string, userId string){
	from := mail.NewEmail("HotSys", "Hotsys@mail.com")
	subject := "Email Verification for HotSys"
	to := mail.NewEmail(username, emailAddress)
	content := mail.NewContent("text/plain", "Click on the link below to verify your email address for " + username + "\n " + utils.ConfirmAdminMailEndpoint+ userId + "\nThis link expires in 7 days.")
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














