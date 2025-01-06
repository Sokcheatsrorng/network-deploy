
```
### List all the function that's chainecode handle Admin 

1.GetallDonation 
2.QueryDonationByDonorName // get by name 
3.GetDonationHistory  // check the history of the donation 
4.DonationExists // Checks specific donation record exists on the ledger

### List all the function that's chainecode handle Peer  

1.GetallDonation // all donation make by peer (client)
2.QueryDonationByDonorName // get by name 
3.GetDonationHistoryByName  // check the history of the donation my peers 
4.CreateDonation  // create donation peers 

```


```
### go start to project with golang 

go mod init your-project-name
go mod tidy

```

```
peer lifecycle chaincode package chaincode-donation.tar.gz --path ./chaincode-donation --label chaincode-donation_v1
peer lifecycle chaincode install chaincode-donation.tar.gz

```

```
// Define The Donation Event Type 
example of the form : 

{
  "eventID": "EVT12345",
  "recipient": "John Doe",
  "eventDescription": "Fundraiser to support underprivileged youth with educational resources",
  "organization": "EduForAll Foundation",
  "timestamp": "2024-11-20T10:30:00Z"
}



```


##

# Donation Chaincode

This project implements a donation management system as a Hyperledger Fabric chaincode. It supports managing donation events and individual donations within a blockchain network. The chaincode enforces business rules, such as ensuring recipients belong to the correct organization (MSP).

---

## Features

### Donation Event Management
- **Create Donation Event**: Adds a new donation event to the ledger.
- **Read Donation Event**: Retrieves the details of a donation event by its ID.
- **Update Donation Event**: Updates the details of an existing donation event.
- **Delete Donation Event**: Deletes a donation event from the ledger.
- **Get All Donation Events**: Fetches all donation events stored on the ledger.

### Donation Management
- **Create Donation**: Registers a new donation for a specific donation event.
- **Read Donation**: Fetches donation details by ID.
- **Update Donation**: Updates the details of an existing donation.
- **Delete Donation**: Deletes a donation entry from the ledger.
- **Get All Donations**: Retrieves all donations from the ledger.
---

## Prerequisites
- cc-tools
- [Hyperledger Fabric v2.x](https://hyperledger-fabric.readthedocs.io/) installed.
- CouchDB configured as the state database.
- `fabric-contract-api-go` for chaincode development in Go.

---

### ==>  API Endpoints

### **Donation Event Endpoints**

1. **CreateDonationEvent**
   - **Description**: Creates a new donation event in the blockchain network.
   - **Parameters**:
     - `id`: Unique identifier for the event.
     - `eventName`: Name of the event.
     - `recipient`: Recipient for the event.
     - `description`: Description of the event.
     - `timestamp`: Timestamp for the event creation.
   - **Business Rules**:
     - Ensures the event does not already exist.
     - Validates that the recipient belongs to the creator's organization (MSP).

2. **ReadDonationEvent**
   - **Description**: Retrieves details of a specific donation event.
   - **Parameters**:
     - `id`: Unique identifier of the event.

3. **UpdateDonationEvent**
   - **Description**: Updates an existing donation event.
   - **Parameters**:
     - `id`: Unique identifier of the event.
     - `recipient`: Updated recipient of the event.
     - `description`: Updated description.
     - `timestamp`: Updated timestamp.

4. **DeleteDonationEvent**
   - **Description**: Deletes a donation event from the ledger.
   - **Parameters**:
     - `id`: Unique identifier of the event.

5. **GetAllDonationEvents**
   - **Description**: Fetches all donation events stored in the ledger.

---

### **Donation Endpoints**

1. **CreateDonation**
   - **Description**: Creates a new donation and associates it with a donation event.
   - **Parameters**:
     - `id`: Unique identifier for the donation.
     - `donationEventID`: ID of the related donation event.
     - `donor`: Donor's name.
     - `amount`: Donation amount.
     - `message`: Message from the donor.
     - `recipient`: Recipient of the donation.

2. **ReadDonation**
   - **Description**: Fetches donation details by its ID.
   - **Parameters**:
     - `id`: Unique identifier of the donation.

3. **UpdateDonation**
   - **Description**: Updates details of an existing donation.
   - **Parameters**:
     - `id`: Unique identifier of the donation.
     - `donationEventID`: Associated event ID.
     - `donor`: Updated donor name.
     - `amount`: Updated donation amount.
     - `message`: Updated donor message.
     - `recipient`: Updated recipient.

4. **DeleteDonation**
   - **Description**: Deletes a donation record from the ledger.
   - **Parameters**:
     - `id`: Unique identifier of the donation.

5. **GetAllDonations**
   - **Description**: Retrieves all donations stored in the ledger.

---

## Initialization

Use the `InitLedger` method to populate the blockchain ledger with some initial test data. This method creates predefined donation events and donations.

---

## Business Rules

- **Organization Validation**: Only recipients belonging to the creator's organization (MSP) are allowed to create donation events.
- **Unique Identifiers**: All donation events and donations must have unique IDs.
- **Timestamp Enforcement**: The current timestamp is used to record the time of creation or update.

---

## Deployment

1. Package the chaincode:
   ```bash
   peer lifecycle chaincode package donation.tar.gz --path ./ --lang golang --label donation_1.0


## test go project 
go test ./... -v 