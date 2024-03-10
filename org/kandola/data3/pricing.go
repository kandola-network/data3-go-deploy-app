package data3

import (
	"encoding/json"
	"log"
)

// calculatePrice takes a DeploymentRequest and returns a price estimation
/*
Sample Request JSON:
{
   "dbOwner": "0xb34ce981c44702f0e5b19884009f76be2f10fbb7",
   "region": "WESTERN_EUROPE-United Kingdom-London",
   "dbEngine": "MySQL",
   "dbEngineVersion": "8.2.0",
   "isLicensed": null,
   "licenseKey": null,
   "deploymentType": "SHARED",
   "specification": "DB STARTER K2 - 1 GB RAM - 2 vCPU(s)",
   "cpu": 2,
   "memory": 1,
   "storage": 100,
   "iops": 1000,
   "isRedundancyRequired": false,
   "redundancy": null,
   "isBackupRequired": false,
   "backupFrequencyDays": null,
   "backupRetentionDays": null,
   "paymentFrequency": "MONTHLY",
   "dbUsername": "admin",
   "isAutoGenPassword": false,
   "dbPassword": "admin",
   "name": "Test 37",
   "address": "c934502630220c58b726194c8ce87e96fb1adf48cf2720a053f8ca7721101238"
}

*/
func calculatePrice(deploymentRequest DeploymentRequest) float64 {
	// This is a sample implementation - replace it with your actual logic

	// Convert the payload back to JSON for logging
	deploymentRequestJSON, err := json.Marshal(deploymentRequest)
	if err != nil {
		log.Printf("Error marshalling pricing request to JSON: %v", err)
	} else {
		// Log the JSON string
		log.Printf("Pricing request received: %s\n", deploymentRequestJSON)
	}

	// Implement Your Pricing logic here...
	// Simple pricing logic
	basePrice := 50.0                                        // Base price for any deployment
	cpuPrice := float64(deploymentRequest.CPU) * 10.0        // Example: $10 per vCPU core
	memoryPrice := float64(deploymentRequest.Memory) * 5.0   // Example: $5 per GB of memory
	storagePrice := float64(deploymentRequest.Storage) * 2.0 // Example: $2 per GB of storage

	redundancyMultiplier := 1.0
	if deploymentRequest.IsRedundancyRequired {
		redundancyMultiplier = float64(*deploymentRequest.
			Redundancy) // To ensure this is correctly parsed from JSON as a float64
	}

	dedicatedMultiplier := 1.0
	if deploymentRequest.DeploymentType == "SHARED" {
		dedicatedMultiplier = 1.0
	} else {
		dedicatedMultiplier = 2.0
	}

	iopsMultiplier := 1.0
	iops := deploymentRequest.IOPS
	switch {
	case iops <= 1000:
		iopsMultiplier = 1.0
	case iops <= 2000:
		iopsMultiplier = 1.2
	case iops <= 3000:
		iopsMultiplier = 1.3
	case iops <= 4000:
		iopsMultiplier = 1.4
	case iops <= 6000:
		iopsMultiplier = 1.6
	case iops <= 9000:
		iopsMultiplier = 1.9
	case iops <= 12000:
		iopsMultiplier = 2.5
	default:
		iopsMultiplier = 10.0
	}

	totalPrice := (basePrice + cpuPrice + memoryPrice + storagePrice) * redundancyMultiplier * dedicatedMultiplier * iopsMultiplier
	return totalPrice
}
