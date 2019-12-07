package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

//Collection parameters
const DatabaseName = "HotelManagementSystemDatabase"
const HotelCollection = "hotels"
const EnvironmentVariableFilename = "EnvironmentVariables.env"

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