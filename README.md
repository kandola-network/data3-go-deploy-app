# data3-go-deploy-app
A GoLang based REST App for automating Data3 Deployments and Pricing

# Sample Data3.Network Pricing and Deployment Engine App (Go)

This Go project is designed to serve as an automated pricing and deployment engine for node providers within the 
Data3 Network's decentralized Database-as-a-Service (DBaaS) platform. 
Utilizing REST APIs, the application receives database deployment requests from customers, 
processes these requests using predefined pricing logic, and relays the calculated prices back to the 
network via a configurable webhook. This system allows node providers to dynamically and efficiently set 
prices and deploy for database services in a decentralized ecosystem.

# Getting Started

## Prerequisites
* Go 1.19+

Check your Go installation by running:
```bash
go version
```
## Clone, Build, Configure, Run
### Clone the repo
```bash
git clone [repository URL]
cd [local repository]
```

### Build using
```bash
go build -o data3-go-deploy-app
```

### Configure using the .env file
```bash
# .env file
WEBHOOK_URL=http://192.168.5.168:8080/internalCommunication/responseForProposal
PORT=3080
HOSTNAME=0.0.0.0
```

### Run
```bash
./data3-go-deploy-app
```

#API

The application exposes the following endpoint for processing pricing requests:
## Pricing
* HTTP Method: POST
* Endpoint: `/price`
* Request Body (JSON example):
```json
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
```
* Calls WEBHOOK_URL to let Data3 Node Provider know of the pricing that was calculated for this request

## Deployment
* HTTP Method: POST
* Endpoint: `/deploy`
* Request Body (JSON example):
```json
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
```

# Built With
* [Go](https://go.dev/)

