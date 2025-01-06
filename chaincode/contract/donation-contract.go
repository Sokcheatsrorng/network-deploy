package contract

import (
	"chaincode-donation/models"
	// "chaincode-donation/txdefs"
	"chaincode-donation/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DonationContract struct {
	contractapi.Contract
}
// text for reponse 
type TxResponse struct {
	TxID    string `json:"txID"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ==============================================================
// -----> Create event inside the Organization // default OrgMSP1
// ==============================================================
// -----> create DonationEvent old response 

// func (dc *DonationContract) CreateDonationEvent(ctx contractapi.TransactionContextInterface, id,eventName,recipient, description, timestamp string) error {
    
//     // Validate event existence
//     exists, err := dc.DonationEventExists(ctx, id)

//     if err != nil {
//         return fmt.Errorf("error checking event existence: %v", err)
//     }
//     if exists {
//         return fmt.Errorf("event with ID %s already exists", id)
//     }

//     // Fetch creator's MSP ID form the api and network 
//     creatorMSPID, err := utils.GetCreatorMSPID(ctx)

//     if err != nil {
//         return fmt.Errorf("error getting creator MSP ID: %v", err)
//     }

//     // Vaildate the Orgnazition with Org vs MSP client 
//     if err := utils.ValidateOrg(ctx,creatorMSPID); err != nil {
        
//        return fmt.Errorf(" Organazition validation failed : %v ",err)

//     }
//     // get the current timesteamp form transaction  
//     currentTimeStampProto , err := ctx.GetStub().GetTxTimestamp()
//     if err != nil {
//         return fmt.Errorf("Faild to get the currentTimeStamp: %v",err)
//     }

//     // the the brotobuf to timetime to time.time
//     currentTimeStamp := time.Unix(currentTimeStampProto.Seconds,int64(currentTimeStampProto.Nanos))
//     // // Format as RFC3339 string
//     currentTimeStampFormat := currentTimeStamp.Format(time.RFC3339)


//     event := models.DonationEvent{
//         ID: id,
// 		EventName: eventName,
//         Recipient:recipient,
//         Description: description,
//         Timestamp:   currentTimeStampFormat,
//         Organization: creatorMSPID,
//         Donations: []models.Donation{},
//     }

//     eventJSON, err := json.Marshal(event)

//     if err != nil {
//         return fmt.Errorf("error marshaling event: %v", err)
//     }

//     return ctx.GetStub().PutState(id, eventJSON)
// }

// CreateDonationEvent handles the creation of a donation event.
func (dc *DonationContract) CreateDonationEvent(ctx contractapi.TransactionContextInterface, id, eventName, recipient, description, timestamp string) (*TxResponse, error) {
	// Ensure the MSP of the creator is valid
	creatorMSP, err := utils.GetCreatorMSPID(ctx)
	if err != nil {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Failed to retrieve creator MSP."}, err
	}

	if err := utils.ValidateOrg(ctx, creatorMSP); err != nil {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Unauthorized MSP."}, err
	}

	// Get the current transaction timestamp
	currentTimeStampProto, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Failed to get transaction timestamp."}, fmt.Errorf("failed to get the current timestamp: %v", err)
	}

	// Convert protobuf timestamp to time.Time
	currentTimeStamp := time.Unix(currentTimeStampProto.Seconds, int64(currentTimeStampProto.Nanos))
	// Format as RFC3339 string
	currentTimeStampFormat := currentTimeStamp.Format(time.RFC3339)

	// Check if the event already exists
	exists, err := dc.DonationEventExists(ctx, id)
	if err != nil {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Error checking event existence."}, err
	}
	if exists {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Donation event already exists."}, fmt.Errorf("donation event %s already exists", id)
	}

	// Create the donation event
	event := models.DonationEvent{
		ID:           id,
		EventName:    eventName,
		Recipient:    recipient,
		Description:  description,
		Timestamp:    currentTimeStampFormat,
		Organization: creatorMSP,
		Donations:    []models.Donation{},
	}

    //  log which organizaton create the donation event 
    fmt.Printf("Donation event %s was created by Organization %s\n",id,event.Organization)

	// Serialize and save the donation event
	eventBytes, err := utils.ToJSON(event)

	if err != nil {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Failed to serialize event."}, err
	}

	err = ctx.GetStub().PutState(id, eventBytes)
    
	if err != nil {
		return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: false, Message: "Failed to save event to ledger."}, err
	}

	return &TxResponse{TxID: ctx.GetStub().GetTxID(), Success: true, Message: "Donation event created successfully."}, nil
}



// --------> Read only DonationEvent by govent ID 
func (dc *DonationContract) ReadDonationEvent(ctx contractapi.TransactionContextInterface, id string) (*models.DonationEvent, error) {
    
    eventJSON, err := ctx.GetStub().GetState(id)
    
    if err != nil {
        return nil, fmt.Errorf("error reading event: %v", err)
    }
    if eventJSON == nil {
        return nil, fmt.Errorf("event with ID %s does not exist", id)
    }

    var event models.DonationEvent

    err = json.Unmarshal(eventJSON, &event)
    if err != nil {
        return nil, fmt.Errorf("error unmarshaling event: %v", err)
    }

    return &event, nil
}

// -------->   UpdateDonation 
func (dc *DonationContract) UpdateDonationEvent(ctx contractapi.TransactionContextInterface, id, recipient, description, timestamp string) error {
    
    // Check if event exists
    exists, err := dc.DonationExists(ctx, id)
    if err != nil {
        return fmt.Errorf("error checking event existence: %v", err)
    }
    if !exists {
        return fmt.Errorf("event with ID %s does not exist", id)
    }

    // get the current timesteamp form transaction  
    currentTimeStampProto , err := ctx.GetStub().GetTxTimestamp()
    if err != nil {
        return fmt.Errorf("Faild to get the currentTimeStamp: %v",err)
    }

    // the the brotobuf to timetime to time.time
    currentTimeStamp := time.Unix(currentTimeStampProto.Seconds,int64(currentTimeStampProto.Nanos))
    // // Format as RFC3339 string
    currentTimeStampFormat := currentTimeStamp.Format(time.RFC3339)

    

    // Update the event
    updatedEvent := models.DonationEvent{
        ID:          id,
        Recipient:   recipient,
        Description: description,
        Timestamp:   currentTimeStampFormat,
        Organization: "", // Retain the organization field or fetch dynamically
    }

    eventJSON, err := json.Marshal(updatedEvent)
    
    if err != nil {
        return fmt.Errorf("error marshaling updated event: %v", err)
    }

    return ctx.GetStub().PutState(id, eventJSON)
}

// -------->  DeleteDonationEvent
func (dc *DonationContract) DeleteDonationEvent(ctx contractapi.TransactionContextInterface, id string) error {
    // Check if event exists
    exists, err := dc.DonationExists(ctx, id)
    if err != nil {
        return fmt.Errorf("error checking event existence: %v", err)
    }
    if !exists {
        return fmt.Errorf("event with ID %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
}

// -------->  GetAllDonationEvents
func (dc *DonationContract) GetAllDonationEvents(ctx contractapi.TransactionContextInterface) ([]*models.DonationEvent, error) {
    
    iterator, err := ctx.GetStub().GetStateByRange("", "")
    
    if err != nil {
        return nil, fmt.Errorf("error retrieving events: %v", err)
    }

    defer iterator.Close()

    var events []*models.DonationEvent

    for iterator.HasNext() {
        response, err := iterator.Next()
        if err != nil {
            return nil, fmt.Errorf("error iterating over events: %v", err)
        }

        var event models.DonationEvent

        err = json.Unmarshal(response.Value, &event)
        if err != nil {
            return nil, fmt.Errorf("error unmarshaling event: %v", err)
        }

        // ensure when it's empty 
        event.InitializeDefaults()

        events = append(events, &event)
    }

    return events, nil
}

// ---------> Check the DonatinEvent Exsit or Not 
func (dc *DonationContract) DonationEventExists(ctx contractapi.TransactionContextInterface,donationEventID string)(bool,error){
   // event 
   eventJSON, err := ctx.GetStub().GetState(donationEventID);

   if err != nil {
	return false,fmt.Errorf("error checking event exsitence: %v",err);
   }
   return eventJSON != nil , nil;
}



// ================================================================
// InitLedger adds a base set of donations to the ledger
// InitLedger just for testing 

func (dc *DonationContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	donations := []models.Donation{
		{ID: "donation1", Donor: "Alice", Amount: 100, Message: "For education", Recipient: "CharityA", Timestamp: "2023-10-10T10:00:00Z"},
		{ID: "donation2", Donor: "Bob", Amount: 150, Message: "Medical aid", Recipient: "CharityB", Timestamp: "2023-10-11T11:30:00Z"},
		{ID: "donation3", Donor: "Charlie", Amount: 200, Message: "Food and shelter", Recipient: "CharityC", Timestamp: "2023-10-12T14:15:00Z"},
		{ID: "donation4", Donor: "Dana", Amount: 250, Message: "Water sanitation", Recipient: "CharityD", Timestamp: "2023-10-13T16:45:00Z"},
		{ID: "donation5", Donor: "Eve", Amount: 300, Message: "Emergency relief", Recipient: "CharityE", Timestamp: "2023-10-14T18:00:00Z"},
	}

	for _, donation := range donations {
		// donationJSON, err := json.Marshal(donation)
        donationJSON, err := json.Marshal(donation)
	    // donationJSON, err := utils.ToJSON(donation)
		if err != nil {
			return fmt.Errorf("failed to marshal donation: %v", err)
		}

		err = ctx.GetStub().PutState(donation.ID, donationJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}


// CreateDonation registers a new donation in CouchDB
func (dc *DonationContract) CreateDonation(
    ctx contractapi.TransactionContextInterface,
    id,donationEventID string,
    donor string,
    amount float32,
    message string,
    recipient string) error {

    // Get the current timestamp
    // donationTimestamp := time.Now().UTC().Format(time.RFC3339)
    // generate a unquine id for a donation 
    id  = utils.GenereteHashID(donor,donationEventID)
    
    // Check if the donation ID already exists
    // exists, err := dc.DonationExists(ctx, id)
    // if err != nil {
    //     return err
    // }
    // if exists {
    //     return fmt.Errorf("Donation with ID %s already exists", id)
    // }
    // go to get the donationEvent 

    eventJSON, err := ctx.GetStub().GetState(donationEventID)
    if err != nil {
        return fmt.Errorf("Failed to retrieve the donationEvent : %v",eventJSON)
    }

    if eventJSON == nil {
        return fmt.Errorf("The DonationEvent with given %s dosen't exsit ",donationEventID)
    }

    // Check if the specified DonationEvent exists
    eventExists, err := dc.DonationEventExists(ctx, donationEventID)
    if err != nil {
        return err
    }
    if !eventExists {
        return fmt.Errorf("donation event with ID %s does not exist", donationEventID)
    }

    // Unmashal  the donationEvent 
    var donationEvent models.DonationEvent
    
    if err := json.Unmarshal(eventJSON, &donationEvent); err != nil {
        return fmt.Errorf("Failed to Unmashal the DonationEvent: %v ",err)
    }

    // get the current timestemp 
    // get the current timesteamp form transaction  
    currentTimeStampProto , err := ctx.GetStub().GetTxTimestamp()
    if err != nil {
        return fmt.Errorf("Faild to get the currentTimeStamp: %v",err)
    }

    // the the brotobuf to timetime to time.time
    currentTimeStamp := time.Unix(currentTimeStampProto.Seconds,int64(currentTimeStampProto.Nanos))
    // // Format as RFC3339 string
    currentTimeStampFormat := currentTimeStamp.Format(time.RFC3339)

    // Create the donation object
    donation := models.Donation{
        ID:            id,
		DonationEventID: donationEventID,
        Donor:         donor,
        Amount:        amount,
        Message:       message,
        Recipient:     recipient,
        Timestamp:     currentTimeStampFormat,
    }

    // Serialize the donation object and save it to the ledger
    donationJSON, err := json.Marshal(donation)

    if err != nil {
        return fmt.Errorf("failed to Marshal the Donation:  %v",err)
    }
   if err := ctx.GetStub().PutState(id,donationJSON); err != nil{
    return fmt.Errorf("Error to create the Donation: %v",err)
   }

   // append new donation to the doationEvent 
   donationEvent.Donations = append(donationEvent.Donations,donation)

    // Serialize and update the DonationEvent in the ledger
    updatedEventJSON, err := json.Marshal(donationEvent)

    if err != nil {
        return fmt.Errorf("failed to marshal updated donation event: %v", err)
    }
    if err := ctx.GetStub().PutState(donationEventID, updatedEventJSON); err != nil {
        return fmt.Errorf("failed to update donation event: %v", err)
    }

    return nil
}


// ReadDonation retrieves a donation from CouchDB by ID
func (dc *DonationContract) ReadDonation(ctx contractapi.TransactionContextInterface, id string) (*models.Donation, error) {
	donationJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read donation: %v", err)
	}
	if donationJSON == nil {
		return nil, fmt.Errorf("donation %s does not exist", id)
	}

	var donation models.Donation
	err = json.Unmarshal(donationJSON, &donation)
	if err != nil {
		return nil, err
	}

	return &donation, nil
}

// UpdateDonation allows updating existing donation details
func (dc *DonationContract) UpdateDonation(ctx contractapi.TransactionContextInterface, id,donationEventID string , donor string, amount float32, message, recipient, timestamp string) error {
	
	exists, err := dc.DonationExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("donation %s does not exist", id)
	}

	// Check if the specified DonationEvent exists
	eventExists, err := dc.DonationEventExists(ctx, donationEventID)
	if err != nil {
		return err
	}
	if !eventExists {
		return fmt.Errorf("donation event with ID %s does not exist", donationEventID)
	}

	// // get the current timestamp 
	// donationTimpstemp := time.Now().UTC().Format(time.RFC3339)
    
    // get the current timesteamp form transaction  
    currentTimeStampProto , err := ctx.GetStub().GetTxTimestamp()
    if err != nil {
        return fmt.Errorf("Faild to get the currentTimeStamp: %v",err)
    }

    // the the brotobuf to timetime to time.time
    currentTimeStamp := time.Unix(currentTimeStampProto.Seconds,int64(currentTimeStampProto.Nanos))
    // // Format as RFC3339 string
    currentTimeStampFormat := currentTimeStamp.Format(time.RFC3339)
    

	donation := models.Donation{
		ID:         id,
		DonationEventID: donationEventID,
		Donor:      donor,
		Amount:     amount,
		Message:    message,
		Recipient:  recipient,
		Timestamp:  currentTimeStampFormat,
	}

	donationJSON, err := json.Marshal(donation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, donationJSON)
}

// DeleteDonation removes a donation entry from CouchDB
func (dc *DonationContract) DeleteDonation(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := dc.DonationExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("donation %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// DonationExists checks if a donation exists in CouchDB by ID
func (dc *DonationContract) DonationExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	donationJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return donationJSON != nil, nil
}

// get function to get all Donatin data : 
// GetAllDonations retrieves all donations from CouchDB
func (dc *DonationContract) GetAllDonations(ctx contractapi.TransactionContextInterface) ([]*models.Donation, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve donations: %v", err)
	}
	
	defer resultsIterator.Close()

	var donations []*models.Donation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var donation models.Donation
		err = json.Unmarshal(queryResponse.Value, &donation)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal donation: %v", err)
		}

		donations = append(donations, &donation)
	}

	return donations, nil
}