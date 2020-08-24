package main

import (
	"encoding/json"
	"fmt"
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
		Endpoint:    aws.String("https://arn:aws:dynamodb:us-east-1:123456789012:table"),
		Credentials: credentials.NewStaticCredentials("AKIAI7YQABX3X4C72V2A", "/iQO1Cma9PT8wJcEYD3ApD4YIl5kHPo10pstbX/N", ""),
	}))
	requestBody := request.Body
	jsonBytes := ([]byte)(requestBody)
	userReq := new(Request)
	fmt.Println(request.Body + "そんなことないよ")
	fmt.Println("通ってる")
	if err := json.Unmarshal(jsonBytes, userReq); err != nil {
		fmt.Println(123)
		fmt.Println("[ERROR]", err)
		return response(http.StatusBadRequest, err.Error()), nil
	}
	user := User{
		LoginId: userReq.LoginId,
		Password: createSafetyPass(userReq.Password),
		UserName: userReq.UserName,
	}
	db := dynamo.New(sess)
	if err := db.Table("user").Put(user).Run(); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}

	return response(http.StatusOK, string("OK")), nil
}

type Request struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
	UserName string `json:"user_name"`
}
type User struct {
	LoginId  string `dynamo:"LoginId,hash"`
	Password string `dynamo:"Password"`
	UserName string `dynamo:"UserName"`
}

func createSafetyPass(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
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