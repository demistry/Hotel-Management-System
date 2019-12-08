package responses

import "strconv"

type GenericResponse struct{
	Status bool `json:"status"`
	Message string `json:"message"`
}

func (resp GenericResponse) AsBytes() []byte{
	respString := "{\n\"status\": " + strconv.FormatBool(resp.Status) + ",\n\"message\" : " + "\"" + resp.Message + "\"\n}"
	return []byte(respString)
}