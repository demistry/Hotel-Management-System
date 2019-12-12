package users

import "net/http"

func CreateUser(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
}
