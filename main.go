package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/scottmckendry/beam/acs"
)

func main() {
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
			Subject:   "Test Email from ACS",
			PlainText: "This is a test email sent using Azure Communication Services.",
			Html:      "<h1>Hello from ACS Email</h1><p>This is a test email sent using Azure Communication Services.</p>",
		},
		Attachments: []acs.Attachment{attachment},
	}

	if err := client.SendEmail(request); err != nil {
		log.Fatal(err)
	}
}
