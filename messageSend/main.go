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
	"math/rand"
	"net/http"
	"time"
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
	fmt.Println(user)
	fmt.Println(r)
	if user.IsLoggedIn == false {
		return response(http.StatusBadRequest, "not logged in"), nil
	}
	//roomが存在するか
	var room []Room
	_ = db.Table("room").Scan().
		Filter("'UserName' = ? AND 'RoomName' = ?", user.UserName, r.RoomName).
		All(&room)
	if len(room) == 0 {
		return response(http.StatusBadRequest, "room not found"), nil
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	dateTIme := time.Now().UTC().In(jst).Format(time.RFC3339)
	fmt.Println(dateTIme)
	message := Message{
		Id:       RandString(20),
		RoomName: r.RoomName,
		UserName: user.UserName,
		Text:     r.Message,
		DateTime: dateTIme,
	}
	if err := db.Table("message").Put(message).Run(); err != nil {
		return response(http.StatusBadRequest, err.Error()), nil
	}
	jsonByte, _ := json.Marshal(message)
	return response(http.StatusOK, string(jsonByte)), nil
}

type Request struct {
	RoomName string `json:"room_name"`
	Message string `json:"message"`
}

type Message struct {
	Id string `json:"id"`
	RoomName string `json:"room_name"`
	UserName string `json:"user_name"`
	Text string `json:"text"`
	DateTime string `json:"date_time"`
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

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandString(n int) string {
	var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}