package models

import "go.mongodb.org/mongo-driver/bson/primitive"
//TODO: Add case for multiple hotels later...
type Room struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RoomName string `json:"roomName,omitempty" bson:"roomName,omitempty"`
	RoomNumber int `json:"roomNumber,omitempty" bson:"roomNumber,omitempty"`
	IsRoomAvailable bool `json:"isRoomAvailable,omitempty" bson:"isRoomAvailable,omitempty"`
	RoomCategory string `json:"roomCategory,omitempty" bson:"roomCategory,omitempty"`
	RoomPrice float64 `json:"roomPrice,omitempty" bson:"roomPrice,omitempty"`
	RoomRank int `json:"roomRank,omitempty" bson:"roomRank,omitempty"`
	RoomImageLink string `json:"roomImageLink,omitempty" bson:"roomImageLink,omitempty"`
	RoomReviews []Reviews `json:"roomReviews,omitempty" bson:"roomReviews,omitempty"`
	HotelID string `json:"hotelId" bson:"hotelId"`
	HotelName string `json:"hotelName" bson:"hotelName"`
}

type Reviews struct{
	ReviewMessage string `json:"reviewMessage" bson:"reviewMessage"`
	ReviewRating string `json:"reviewRating" bson:"reviewRating"`
	Reviewer RoomUser `json:"reviewer" bson:"reviewer"`
}

type RoomUser struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName string `json:"lastName" bson:"lastName"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	ImageLink string `json:"imageLink" bson:"imageLink"`
}

type AdminUser struct{
	ID primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	HotelName string `json:"hotelName" bson:"hotelName"`
	HotelAddress string `json:"hotelAddress" bson:"hotelAddress"`
	HotelRank int `json:"hotelRank" bson:"hotelRank"`
	HotelEmail string `json:"hotelEmail" bson:"hotelEmail"`
	HotelPassword string `json:"hotelPassword" bson:"hotelPassword"`
}
