package responses

import (
	"strconv"
	"time"
)

type GenericResponse struct{
	Status bool `json:"status"`
	Message string `json:"message"`
}

type SuccessfulResponse struct{
	Status bool `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type AdminUserResponse struct{
	ID string `json:"id,omitempty" bson:"id,omitempty"`
	HotelName string `json:"hotelName" bson:"hotelName"`
	HotelAddress string `json:"hotelAddress" bson:"hotelAddress"`
	HotelRank int `json:"hotelRank" bson:"hotelRank"`
	HotelEmail string `json:"hotelEmail" bson:"hotelEmail"`
	IsUserVerified bool `json:"isUserVerified,omitempty" bson:"isUserVerified"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
}

type HotelUserResponse struct {
	ID string `json:"id,omitempty" bson:"id,omitempty"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName string `json:"lastName" bson:"lastName"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	ImageLink string `json:"imageLink" bson:"imageLink"`
}

func (resp GenericResponse) AsBytes() []byte{
	respString := "{\n\"status\": " + strconv.FormatBool(resp.Status) + ",\n\"message\" : " + "\"" + resp.Message + "\"\n}"
	return []byte(respString)
}