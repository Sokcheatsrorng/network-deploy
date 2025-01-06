package txdefs

import (
	"chaincode-donation/utils"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateDonationEvent = tx.Transaction{
	Tag:         "updateDonationEvent",
	Label:       "Update Donation Event",
	Description: "Updates an existing donation event in the ledger",
	Method:      "PUT",
	Callers: []accesscontrol.Caller{
		{
			MSP: "org1MSP",
			OU:  "admin",
		},
		{
			MSP: "org2MSP",
			OU:  "admin",
		},
	},
	Args: []tx.Argument{
		{
			Tag:      "eventID",
			Label:    "Event ID",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "eventName",
			Label:    "Event Name",
			DataType: "string",
			Required: false,
		},
		{
			Tag:      "recipient",
			Label:    "Recipient",
			DataType: "string",
			Required: false,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		
		// Build the update payload from request arguments
		updatePayload := make(map[string]interface{})
		updatePayload["@assetType"] = "donationEvent"
		updatePayload["@key"] = req["eventID"]

		// Add fields to be updated
		if name, ok := req["eventName"].(string); ok {
			updatePayload["eventName"] = name
		}
		if recipient, ok := req["recipient"].(string); ok {
			updatePayload["recipient"] = recipient
		}

		// Use the Update method to apply changes
		key, _ := assets.NewKey(updatePayload)
		updatedAsset, err := key.Update(stub, updatePayload)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update donation event asset")
		}

		// Marshal the updated asset into JSON for the response
		response, err := utils.ToJSON(updatedAsset)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to marshal updated asset to JSON")
		}

		return response, nil
	},
}
