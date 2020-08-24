.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o signUp/signUp ./signUp/main.go
	GOOS=linux GOARCH=amd64 go build -o signIn/signIn ./signIn/main.go
	GOOS=linux GOARCH=amd64 go build -o roomCreate/roomCreate ./roomCreate/main.go
	GOOS=linux GOARCH=amd64 go build -o roomGet/roomGet ./roomGet/main.go
	GOOS=linux GOARCH=amd64 go build -o roomInvite/roomInvite ./roomInvite/main.go
	GOOS=linux GOARCH=amd64 go build -o messageSend/messageSend ./messageSend/main.go
	GOOS=linux GOARCH=amd64 go build -o messageGet/messageGet ./messageGet/main.go
	GOOS=linux GOARCH=amd64 go build -o searchUser/searchUser ./searchUser/main.go

package:
	sam	package --template-file template.yaml --output-template-file output-template.yaml --s3-bucket sam-template-store-go-chatt --profile kendoyasui

deploy:
	sam deploy --template-file output-template.yaml --stack-name sam-template-store-go-chatt --capabilities CAPABILITY_IAM --profile kendoyasui