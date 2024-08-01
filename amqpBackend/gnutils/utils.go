package gnutils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(err)
	}

}
func LoadEnv() {
	err := godotenv.Load()
	failOnError(err, "Error loading .env file")
}

// SendMail sends an email using the provided SMTP server details and authentication.
func SendMail(to []string, subject, body string) {
	LoadEnv()
	smtpSvr := os.Getenv("SMTP_SERVER")
	authEmail := os.Getenv("AUTH_EMAIL")
	authPass := os.Getenv("AUTH_PASSWORD")

	auth := smtp.PlainAuth("", authEmail, authPass, smtpSvr)

	// Email details

	from := authEmail
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", smtpSvr+":465", &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpSvr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpSvr)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	// Set the sender and recipients
	if err = client.Mail(from); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			log.Fatalf("Failed to set recipient: %v", err)
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to send email body: %v", err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Fatalf("Failed to write email body: %v", err)
	}
	err = w.Close()
	if err != nil {
		log.Fatalf("Failed to close email body writer: %v", err)
	}

	// Quit the SMTP session
	if err = client.Quit(); err != nil {
		log.Fatalf("Failed to quit SMTP session: %v", err)
	}

	fmt.Println("Email sent successfully!")
}

func PublishMessage(QOE string, typeDeclare string, exchKind string, messageBody map[string]interface{}, routingKey string) {
	LoadEnv()

	rabbitmqURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitmqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	fmt.Println("Success Connected")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ttl := int32(10000)
	message := messageBody
	messageBodyJson, err := json.Marshal(message)
	failOnError(err, "Failed to marshal JSON")

	if typeDeclare == "exchange" {
		err := ch.ExchangeDeclare(
			QOE,
			exchKind,
			true,
			false,
			false,
			false,
			amqp.Table{"x-message-ttl": ttl},
		)
		failOnError(err, "Failed to declare a exchange")
	} else if typeDeclare == "queue" {
		_, err := ch.QueueDeclare(
			QOE,
			true,
			false,
			false,
			false,
			amqp.Table{"x-message-ttl": ttl},
		)
		failOnError(err, "Failed to declare a queue")
	}
	err = ch.Publish(
		QOE,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBodyJson,
		},
	)
	failOnError(err, "Failed to publish a message")

	fmt.Println("Success Publish message to queue")
}

func ConsumeMessage(queue string) <-chan amqp.Delivery {
	LoadEnv()

	rabbitmqURL := os.Getenv("RABBITMQ_URL")

	fmt.Println("Application Consumer")
	conn, err := amqp.Dial(rabbitmqURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	fmt.Println("Success Connected")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	msgs, err := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	fmt.Println("Consumer registered, waiting for messages...")
	return msgs
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
