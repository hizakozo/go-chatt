package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"net/http"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String("http://dynamoTest:8000"),
		Credentials: credentials.NewStaticCredentials("tekitou", "tekitou", ""),
	}))
	requestBody := request.Body
	jsonBytes := ([]byte)(requestBody)
	r := new(Request)
	if err := json.Unmarshal(jsonBytes, r); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	db := dynamo.New(sess)
	var user User
	if err := db.Table("user").Get("LoginId", request.Headers["Login-Id"]).One(&user); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	if user.IsLoggedIn == false {
		return response(http.StatusBadRequest, "not logged in"), nil
	}
	room := Room{
		UserName: user.UserName,
		RoomName: r.RoomName,
	}
	if err := db.Table("room").Put(room).Run(); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	jsonByte, _ := json.Marshal(ResponseBody{RoomName: room.RoomName})
	return response(http.StatusOK, string(jsonByte)), nil
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