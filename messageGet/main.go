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
	"net/http"
	"sort"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		//Endpoint:    aws.String("http://dynamoTest:8000"),
		Credentials: credentials.NewStaticCredentials("AKIAJRYKCFVFXH4VABZA", "mYNWIsEiRXecxT6EJbW0cD06jsk78H0J+pwFy/Mp", ""),
	}))
	db := dynamo.New(sess)
	fmt.Println(request.QueryStringParameters["room_name"])
	roomName := request.QueryStringParameters["room_name"]
	//roomが存在するか
	var room []Room
	_ = db.Table("room").Scan().
		Filter("'UserName' = ? AND 'RoomName' = ?", "yasui", roomName).
		All(&room)
	if len(room) == 0 {
		return response(http.StatusBadRequest, "room not found"), nil
	}
	var messages Messages
	_ = db.Table("message").Scan().
		Filter("'RoomName' = ?", roomName).
		All(&messages)
	sort.Sort(ByDateTime{messages})
	responseJson, _ := json.Marshal(ResponseBody{RoomName: roomName, Messages: messages})
	return response(http.StatusOK, string(responseJson)), nil
}

type ResponseBody struct {
	RoomName string    `json:"room_name"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Id       string `json:"id"`
	RoomName string `json:"room_name,omitempty"`
	UserName string `json:"user_name"`
	Text     string `json:"text"`
	DateTime string `json:"date_time"`
}

type Messages []Message
type ByDateTime struct {
	Messages
}
func (b ByDateTime) Less(i, j int) bool {
	return b.Messages[i].DateTime < b.Messages[j].DateTime
}
func (m Messages) Len() int {
	return len(m)
}
func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type Room struct {
	UserName string `json:"user_name"`
	RoomName string `json:"room_name"`
}

type User struct {
	LoginId    string `dynamo:"LoginId,hash"`
	Password   string `dynamo:"Password"`
	UserName   string `dynamo:"UserName"`
	IsLoggedIn bool   `dynamo:"IsLoggedIn"`
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Headers": "Content-Type,login_id",
			"Access-Control-Allow-Methods": "GET,OPTIONS,POST",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
