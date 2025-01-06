package txdefs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// CreateDonationEvent handles the creation of donation events in the blockchain.
var CreateDonationEvent = tx.Transaction{
	Tag:         "createDonationEvent",
	Label:       "Create Donation Event",
	Description: "Create a new donation event",
	Method:      "POST",
	// who can  call this 
	Callers: []accesscontrol.Caller{
		{
			MSP: "Org1MSP",
			OU: "admin",
		},
		{
			// will  change it only org1 can create donationEvent
			MSP: "Org2MSP",
			OU: "admin",
		},
	},
	Args: []tx.Argument{
		{
			Tag:         "id",
			Label:       "Event ID",
			Description: "Unique identifier for the donation event",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "eventName",
			Label:       "Event Name",
			Description: "Name of the donation event",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "recipient",
			Label:       "Recipient",
			Description: "Recipient organization",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "description",
			Label:       "Description",
			Description: "Description of the donation event",
			DataType:    "string",
			Required:    true,
		},
		{
			
			Tag:         "timestamp",
            Label:       "Timestamp",
            Description: "Event creation timestamp.",
            DataType:    "datetime",
			Required: true,

		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// Extract arguments
		id, _ := req["id"].(string)
		eventName, _ := req["eventName"].(string)
		recipient, _ := req["recipient"].(string)
		description, _ := req["description"].(string)

		// Check if the event already exists
		existingEvent, _ := stub.GetState(id)
		if existingEvent != nil {
			return nil, errors.NewCCError("Donation event already exists", 400)
		}

		// Validate timestamp
		currentTimeProto, err := stub.Stub.GetTxTimestamp()
		if err != nil {
			return nil, errors.WrapError(err, "Failed to retrieve current timestamp")
		}

		// Convert Protobuf timestamp to time.Time
		currentTime := time.Unix(currentTimeProto.Seconds, int64(currentTimeProto.Nanos))
		timestampFormatted := currentTime.Format(time.RFC3339)

		// Create donation event asset
		eventAsset := map[string]interface{}{
			"@assetType":  "donationEvent",
			"id":          id,
			"eventName":   eventName,
			"recipient":   recipient,
			"description": description,
			"timestamp":   timestampFormatted,
			"donations":   []interface{}{},
		}

		// Initialize new asset
		asset, err := assets.NewAsset(eventAsset)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create new donation event asset")
		}

		// Save the new donation event to the ledger
		// var err = error
		_, err = asset.PutNew(stub)

		if err != nil {
			if iccErr, ok := err.(*errors.CCError); ok{

				return nil, errors.WrapErrorWithStatus(iccErr, "Failed to save donation event to ledger", iccErr.Status())
		
			}
		   // In case the error is from cc-errors
		   return nil,errors.WrapError(err,"Failed to save donation event to ledger!") 
		
	    }
			

		// Log event creation
		eventLog := fmt.Sprintf("Donation event created: %s", id)
		events.CallEvent(stub, "DonationEventCreated", []byte(eventLog))

		// Return the created event
		eventJSON, err := json.Marshal(asset)

		if err != nil {
			return nil, errors.WrapError(err, "Failed to serialize donation event")
		}

		// Marshhal message o be logged 
		logMsg, ok := json.Marshal(fmt.Sprintf("New Donation Event naem : %s",eventName))
        if ok != nil {
			return nil,errors.WrapError(nil,"Failed to encode asset to JSON format!")
		}

		// Call event to log the message 
		events.CallEvent(stub,"Create Donation Event Log",logMsg)

		return eventJSON, nil
	},
}
