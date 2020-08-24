package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"golang.org/x/crypto/bcrypt"
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
	userReq := new(Request)
	if err := json.Unmarshal(jsonBytes, userReq); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	db := dynamo.New(sess)
	table := db.Table("user")
	var user User
	if err := table.Get("LoginId", userReq.LoginId).One(&user); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	if err := passwordVerify(user.Password, userReq.Password); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	user.IsLoggedIn = true
	if err := table.Put(user).Run(); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	responseBody, _ := json.Marshal(ResponseBody{LoginId: userReq.LoginId})
	return response(http.StatusOK, string(responseBody)), nil
}

type Request struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

type ResponseBody struct {
	LoginId string `json:"login_id"`
}

func passwordVerify(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
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