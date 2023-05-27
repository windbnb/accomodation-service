package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/util"
)

func GetUserById(userId uint) (model.UserResponseDTO, error) {
	response, err := http.Get(util.BaseUserServicePathRoundRobin.Next().Host + "/api/users/" + fmt.Sprint(userId))
	if err != nil {
		return model.UserResponseDTO{}, err
	}

	var userResponse model.UserResponseDTO
	json.NewDecoder(response.Body).Decode(&userResponse)
	return userResponse, nil
}

func AuthorizeHost(tokenString string) (model.UserResponseDTO, error) {
	url := util.BaseUserServicePathRoundRobin.Next().Host + "/api/users/authorize/host"
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", tokenString)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return model.UserResponseDTO{}, err
	}

	var userResponse model.UserResponseDTO
	json.NewDecoder(response.Body).Decode(&userResponse)
	return userResponse, nil
}
