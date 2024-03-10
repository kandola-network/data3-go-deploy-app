package data3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type DeploymentRequest struct {
	DbOwner              string  `json:"dbOwner"`
	Region               string  `json:"region"`
	DbEngine             string  `json:"dbEngine"`
	DbEngineVersion      string  `json:"dbEngineVersion"`
	IsLicensed           *bool   `json:"isLicensed"` // Using pointer to bool for nullable fields
	LicenseKey           *string `json:"licenseKey"`
	DeploymentType       string  `json:"deploymentType"`
	Specification        string  `json:"specification"`
	CPU                  int     `json:"cpu"`
	Memory               int     `json:"memory"`
	Storage              int     `json:"storage"`
	IOPS                 int     `json:"iops"`
	IsRedundancyRequired bool    `json:"isRedundancyRequired"`
	Redundancy           *int    `json:"redundancy"`
	IsBackupRequired     bool    `json:"isBackupRequired"`
	BackupFrequencyDays  *int    `json:"backupFrequencyDays"`
	BackupRetentionDays  *int    `json:"backupRetentionDays"`
	PaymentFrequency     string  `json:"paymentFrequency"`
	DbUsername           string  `json:"dbUsername"`
	IsAutoGenPassword    bool    `json:"isAutoGenPassword"`
	DbPassword           *string `json:"dbPassword"`
	Name                 string  `json:"name"`
	Address              string  `json:"address"`
}

// priceHandler handles the POST requests to the /price endpoint
func priceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed for /price", http.StatusMethodNotAllowed)
		return
	}

	var payload DeploymentRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the payload for pricing...
	price := calculatePrice(payload)

	if err := postToWebhook(price, payload.Address); err != nil {
		log.Printf("Error posting to webhook on node provider service: %v\n", err)
		http.Error(w, "Failed to call webhook on node provider service", http.StatusInternalServerError)
		return
	}

	successfulMessage := fmt.Sprintf("Sent price of %f for request %s to %s\n",
		price, payload.Address, os.Getenv("WEBHOOK_URL"))
	log.Printf(successfulMessage)
	// Send response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": successfulMessage})
}

// deployHandler handles the POST requests to the /deploy endpoint
func deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed for /deploy", http.StatusMethodNotAllowed)
		return
	}

	var payload DeploymentRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the payload for deployment...
	go deploy(payload)

	successfulMessage := fmt.Sprintf("Deployment request received for request %s\n", payload.Address)
	log.Printf(successfulMessage)
	// Send response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": successfulMessage})
}

func createFormData(price float64, address string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if err := writer.WriteField("rspAddress", address); err != nil {
		return nil, "", err
	}
	if err := writer.WriteField("price", fmt.Sprintf("%f", price)); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func postToWebhook(price float64, address string) error {
	formData, contentType, err := createFormData(price, address)
	if err != nil {
		return err
	}

	webhookURL := os.Getenv("WEBHOOK_URL") // Read the webhook URL from an environment variable
	if webhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL is not set")
	}

	resp, err := http.Post(webhookURL, contentType, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error when invoking webhook %s: HTTP %d %v", webhookURL, resp.StatusCode, resp.Status)
	}

	return nil
}

func Serve() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Reading environment variables for port, hostname, and webhook URL
	port := os.Getenv("PORT")
	if port == "" {
		port = "3080" // Default to 3080 if no environment variable is set
	}

	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = "0.0.0.0" // Default to 0.0.0.0 if no environment variable is set
	}

	webhookURL := os.Getenv("WEBHOOK_URL") // This may remain unset

	// Setting up HTTP server routes
	http.HandleFunc("/price", priceHandler)
	http.HandleFunc("/deploy", deployHandler)

	// Construct the address and start listening
	address := fmt.Sprintf("%s:%s", hostname, port)

	// Log the startup messages
	log.Printf("Pricing / Deployment Engine is running on http://%s\n", address)
	if webhookURL != "" {
		log.Printf("Pricing / Deployment Engine will send calculated prices to %s\n", webhookURL)
	} else {
		log.Println("WEBHOOK_URL environment variable is not set")
		log.Fatalf("Failed to start the server: %v\n", fmt.Errorf("WEBHOOK_URL environment variable is not set"))
	}

	// Start the HTTP server
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start the server: %v\n", err)
	}
}
