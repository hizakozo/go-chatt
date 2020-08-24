#!/bin/bash

aws dynamodb --profile local --endpoint-url http://localhost:8000 create-table --cli-input-json file://./go-chatt/data/user.json
aws dynamodb --profile local --endpoint-url http://localhost:8000 create-table --cli-input-json file://./go-chatt/data/message.json
aws dynamodb --profile local --endpoint-url http://localhost:8000 create-table --cli-input-json file://./go-chatt/data/room.json