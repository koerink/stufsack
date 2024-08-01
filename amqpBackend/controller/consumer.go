package main

import (
	"amqpBackend/gnutils"
	"amqpBackend/models"
	"amqpBackend/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(err)
	}
}

func main() {
	msgs := gnutils.ConsumeMessage("trxRegQueue")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var message map[string]interface{}
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				fmt.Printf("Failed to unmarshal JSON: %s\n", err)
				continue
			}

			fmt.Printf("Received Message: %v\n", message["Type"])

			bodyMap, ok := message["Body"].(map[string]interface{})
			if !ok {
				fmt.Println("Invalid message body")
				continue
			}

			// Extract fields from the bodyMap
			username, _ := bodyMap["username"].(string)
			email, _ := bodyMap["Email"].(string)
			password, _ := bodyMap["Password"].(string)
			isVerifyEmail, _ := bodyMap["IsVerifyEmail"].(string) // Convert to bool later
			profileMap, _ := bodyMap["Profile"].(map[string]interface{})
			courseIds, _ := bodyMap["CourseId"].([]interface{})

			// Convert profile to []string
			profile := make([]string, 0, len(profileMap))
			for _, v := range profileMap {
				profile = append(profile, v.(string))
			}

			// Convert courseIds to []int
			courseIdInts := make([]int, 0, len(courseIds))
			for _, v := range courseIds {
				courseIdInts = append(courseIdInts, int(v.(float64)))
			}

			// Convert isVerifyEmail to bool
			isVerifyEmailBool := isVerifyEmail == "true"

			// Initialize MongoDB client
			mongoTestClient := gnutils.ClientMongo()

			defer func() {
				if err := mongoTestClient.Disconnect(context.Background()); err != nil {
					log.Fatal("Error while disconnecting from DB", err)
				}
			}()

			coll := mongoTestClient.Database("backendAPI").Collection("account_coll")
			accountColl := services.AccountModel{MongoCollection: coll}

			passwordHash, err := gnutils.HashPassword(password)
			if err != nil {
				fmt.Println("Error hashing password", err)
				continue
			}

			account := models.CollAccount{
				Username:      username,
				Email:         email,
				Password:      passwordHash,
				IsVerifyEmail: isVerifyEmailBool,
				Profile:       profile,
				CourseId:      courseIdInts,
			}

			result, err := accountColl.CreateAccount(&account)
			if err != nil {
				fmt.Println("Insert 1 operation failed", err)
				continue
			}

			fmt.Println("Insert 1 Success", result)

			// Prepare and send email
			to := []string{"derry.nur.syariffudin@gmail.com"}
			subject := "Test Email"

			messageType, _ := message["Type"].(string)
			bodyContentBytes, _ := json.MarshalIndent(bodyMap, "", "  ")
			bodyContent := string(bodyContentBytes)

			// Format the email body
			body := fmt.Sprintf("This is a test email sent from a Go application.\nType: %s\nBody: %s", messageType, bodyContent)
			gnutils.SendMail(to, subject, body)
		}
	}()

	fmt.Println("Success receive message")
	fmt.Println("[*] - Wait for message")
	<-forever
}
