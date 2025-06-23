// Package acs provides a client for Azure Communication Services email functionality
package acs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const apiVersion = "2023-03-31"

// EmailClient represents an Azure Communication Services email client.
// It handles authentication and communication with the Azure Communication Services email API.
type EmailClient struct {
	endpoint   string
	credential string
	httpClient *http.Client
}

// EmailRequest represents the structure for sending an email.
type EmailRequest struct {
	SenderAddress string         `json:"senderAddress"`
	Content       EmailContent   `json:"content"`
	Recipients    Recipients     `json:"recipients"`
	Attachments   []Attachment   `json:"attachments,omitempty"`
	ReplyTo       []EmailAddress `json:"replyTo,omitempty"`
}

// Recipients represents a collection of email recipients categorized by type (To, Cc, Bcc).
type Recipients struct {
	To  []EmailAddress `json:"to"`
	Cc  []EmailAddress `json:"cc,omitempty"`
	Bcc []EmailAddress `json:"bcc,omitempty"`
}

// EmailAddress represents an email address with an optional display name.
type EmailAddress struct {
	Address     string `json:"address"`
	DisplayName string `json:"displayName,omitempty"`
}

// EmailContent represents the content of an email including subject and body.
type EmailContent struct {
	Subject   string `json:"subject"`
	PlainText string `json:"plainText,omitempty"`
	Html      string `json:"html,omitempty"`
}

// Attachment represents an email attachment with its metadata and content.
type Attachment struct {
	Name            string `json:"name"`
	ContentType     string `json:"contentType"`
	ContentInBase64 string `json:"contentInBase64"`
}

// EmailOperationStatus represents the status of an email send operation
type EmailOperationStatus struct {
	Status string `json:"status"`
	Error  struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// NewEmailClient creates a new ACS email client with the specified endpoint and credential.
func NewEmailClient(endpoint, credential string) *EmailClient {
	return &EmailClient{
		endpoint:   endpoint,
		credential: credential,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SendEmail sends an email using the Azure Communication Services email API.
func (c *EmailClient) SendEmail(emailReq EmailRequest) (string, error) {
	url := fmt.Sprintf("%s/emails:send?api-version=%s", c.endpoint, apiVersion)

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal email request: %w", err)
	}

	contentHash := generateContentHash(jsonData)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if err := c.signRequest(req, contentHash); err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	operationLocation := resp.Header.Get("Operation-Location")
	return operationLocation, nil
}

func (c *EmailClient) GetEmailStatus(operationLocation string) (*EmailOperationStatus, error) {
	contentHash := generateContentHash([]byte{})

	req, err := http.NewRequest(http.MethodGet, operationLocation, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := c.signRequest(req, contentHash); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var status EmailOperationStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

// NewAttachment creates a new email attachment from the provided file data.
func NewAttachment(name string, contentType string, data []byte) Attachment {
	return Attachment{
		Name:            name,
		ContentType:     contentType,
		ContentInBase64: base64.StdEncoding.EncodeToString(data),
	}
}

// generateContentHash creates a SHA256 hash of the content and returns it as a base64 encoded string.
func generateContentHash(content []byte) string {
	hash := sha256.Sum256(content)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// signRequest adds the necessary authentication headers to the request.
func (c *EmailClient) signRequest(req *http.Request, contentHash string) error {
	timestamp := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	host := strings.TrimPrefix(c.endpoint, "https://")

	verb := strings.ToUpper(req.Method)
	uriPathAndQuery := req.URL.RequestURI()

	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s",
		verb,
		uriPathAndQuery,
		timestamp,
		host,
		contentHash,
	)

	key, err := base64.StdEncoding.DecodeString(c.credential)
	if err != nil {
		return fmt.Errorf("failed to decode access key: %w", err)
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req.Header.Set("x-ms-date", timestamp)
	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("Host", host)
	req.Header.Set("repeatability-request-id", uuid.New().String())
	req.Header.Set("repeatability-first-sent", time.Now().UTC().Format(time.RFC3339))

	authHeader := fmt.Sprintf("HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature=%s", signature)
	req.Header.Set("Authorization", authHeader)

	return nil
}
