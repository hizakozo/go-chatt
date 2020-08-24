package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestBody := request.Body
	jsonBytes := ([]byte)(requestBody)
	r := new(Request)
	if err := json.Unmarshal(jsonBytes, r); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	return response(http.StatusOK, string(jsonBytes)), nil
}

type Request struct {
	RoomName string `json:"room_name"`
}

type ResponseBody struct {
	RoomName string `json:"room_name"`
}

type Room struct {
	UserName  string `json:"user_name"`
	RoomName string `json:"room_name"`
}

type User struct {
	LoginId  string `dynamo:"LoginId,hash"`
	Password string `dynamo:"Password"`
	UserName string `dynamo:"UserName"`
	IsLoggedIn bool `dynamo:"IsLoggedIn"`
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{
			"Access-Control-Allow-Headers": "Content-Type,login_id",
			"Access-Control-Allow-Methods": "GET,OPTIONS,POST",
			"Access-Control-Allow-Origin": "*",
		},
	}
}