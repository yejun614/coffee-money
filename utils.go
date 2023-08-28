package main

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/hex"
	"encoding/json"
	"crypto/sha512"
)

func SHA512(input string) string {
	hash := sha512.New()
	encryped := hash.Sum([]byte(input))
	return hex.EncodeToString(encryped)
}

func GetGithubAccessToken(clientID string, clientSecrets string, code string) (*ResponseAccessTokenGithubDTO, error) {
	// result
	result := new(ResponseAccessTokenGithubDTO)

	// data
	data, err := json.Marshal(&RequestAccessTokenGithubDTO{
		ClientID: clientID,
		ClientSecrets: clientSecrets,
		Code: code,
	})
	if err != nil {
		return result, err
	}

	// make request
	request, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return result, err
	}

	// set request header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// send request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return result, err
	}

	// check response
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	// json parsing
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	// Done
	return result, nil
}

func GetGithubAuthUser(accessToken string) (*ResponseGetAuthenticatedUserGithubDTO, error) {
	// result
	result := new(ResponseGetAuthenticatedUserGithubDTO)

	// make request
	request, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return result, err
	}

	// set request header
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// send request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return result, err
	}

	// check response
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	// json parsing
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	// Done
	return result, nil
}
