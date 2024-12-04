package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	aisstream "github.com/aisstream/ais-message-models/golang/aisStream"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to WebSocket stream
	url := "wss://stream.aisstream.io/v0/stream"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer ws.Close()

	// Subscribe to the WebSocket feed
	subMsg := aisstream.SubscriptionMessage{
		APIKey:          os.Getenv("SECRET_KEY"),
		BoundingBoxes:   [][][]float64{{{49.0, -0.1}, {54.0, 8.0}}}, // bounding box for the entire world
		FiltersShipMMSI: []string{"244592000", "244057043", "244860146", "244860146", "245894000", "244810627", "244083000", "244100829"},
	}

	subMsgBytes, _ := json.Marshal(subMsg)
	if err := ws.WriteMessage(websocket.TextMessage, subMsgBytes); err != nil {
		log.Fatalln(err)
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1")) // Use the appropriate region for your DynamoDB
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create a DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(cfg)

	for {
		// Read incoming WebSocket messages
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}

		var packet aisstream.AisStreamMessage
		err = json.Unmarshal(p, &packet)
		if err != nil {
			log.Fatalln(err)
		}

		var shipName string
		// field may or may not be populated
		if packetShipName, ok := packet.MetaData["ShipName"]; ok {
			shipName = packetShipName.(string)
		}

		// Handle position reports
		switch packet.MessageType {
		case aisstream.POSITION_REPORT:
			var positionReport aisstream.PositionReport
			positionReport = *packet.Message.PositionReport

			// Print the position report
			fmt.Printf("MMSI: %d Ship Name: %s Latitude: %f Longitude: %f\n",
				positionReport.UserID, shipName, positionReport.Latitude, positionReport.Longitude)

			// Prepare item for DynamoDB
			item := map[string]types.AttributeValue{
				"MMSI":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", positionReport.UserID)},
				"Name":      &types.AttributeValueMemberS{Value: shipName},
				"Latitude":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", positionReport.Latitude)},
				"Longitude": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", positionReport.Longitude)},
				"TS":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
			}

			// Push the position data to DynamoDB
			_, err := dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
				TableName: aws.String("ShipCords"),
				Item:      item,
			})
			if err != nil {
				var ae awserr.RequestFailure
				if ok := errors.As(err, &ae); ok {
					log.Printf("Request failed: %v", ae)
				} else {
					log.Fatalf("Failed to insert item into DynamoDB: %v", err)
				}
			} else {
				log.Printf("Successfully inserted position for MMSI: %d", positionReport.UserID)
			}
		}
	}
}
