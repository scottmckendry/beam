package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/scottmckendry/beam/acs"
	"github.com/scottmckendry/beam/db"
)

const skipSend = true

func main() {
	// Initialise the database
	if err := db.InitialiseDb(); err != nil {
		log.Fatal("Failed to initialise database:", err)
		return
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	// Get environment variables
	endpoint := os.Getenv("ACS_ENDPOINT")
	if endpoint == "" {
		log.Fatal("ACS_ENDPOINT is required")
	}

	accessKey := os.Getenv("ACS_ACCESS_KEY")
	if accessKey == "" {
		log.Fatal("ACS_ACCESS_KEY is required")
	}

	senderAddress := os.Getenv("SENDER_ADDRESS")
	if senderAddress == "" {
		log.Fatal("SENDER_ADDRESS is required")
	}

	if skipSend {
		return
	}

	client := acs.NewEmailClient(endpoint, accessKey)

	// Create an attachment
	pdfData, err := os.ReadFile("example.pdf")
	if err != nil {
		log.Fatal("Failed to read PDF file:", err)
	}
	attachment := acs.NewAttachment(
		"report.pdf",
		"application/pdf",
		pdfData,
	)

	request := acs.EmailRequest{
		SenderAddress: senderAddress,
		Recipients: acs.Recipients{
			To: []acs.EmailAddress{
				{
					Address:     "example@example.com",
					DisplayName: "Scott McKendry",
				},
			},
		},
		Content: acs.EmailContent{
			// TODO: all emails should contain opt-out links that call back to the application, updating a db flag in contacts table
			Subject:   "Test Email from ACS",
			PlainText: "This is a test email sent using Azure Communication Services.",
			Html:      "<h1>Hello from ACS Email</h1><p>This is a test email sent using Azure Communication Services.</p>",
		},
		Attachments: []acs.Attachment{attachment},
	}

	operationLocation, err := client.SendEmail(request)

	if err != nil {
		log.Fatal(err)
	}

	// Poll for status up to 30 times with 1 second interval
	for range 30 {
		status, err := client.GetEmailStatus(operationLocation)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Status: %s", status.Status)
		if status.Error.Code != "" {
			// TODO: gracefully handle bounces here, important for domain reputation
			log.Printf("ErrorCode: %s", status.Error.Code)
			log.Printf("ErrorMessage: %s", status.Error.Message)
			break
		}

		if status.Status == "Succeeded" {
			break
		}

		time.Sleep(time.Second)
	}
}
