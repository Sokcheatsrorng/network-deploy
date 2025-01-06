package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// CreateDonation creates a new donation record.
var CreateDonation = tx.Transaction{
	Tag:         "createDonation",
	Label:       "Create Donation",
	Description: "Create a new donation record",
	Method:      "POST",
	Callers: []accesscontrol.Caller{
		{
			MSP: "OrgMSP1",
		},
		{
			MSP: "OrgMSP2",
		},
	},
	Args: []tx.Argument{
		{
			Tag:         "id",
			Label:       "Donation ID",
			Description: "Unique identifier for the donation",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "donor",
			Label:       "Donor",
			Description: "Donor's name or organization",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "amount",
			Label:       "Amount",
			Description: "Amount of the donation",
			DataType:    "float",
			Required:    true,
		},
		{
			Tag:         "eventId",
			Label:       "Event ID",
			Description: "ID of the related donation event",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// Create a new donation asset
		newDonation := map[string]interface{}{
			"@assetType": "donation",
			"id":         req["id"],
			"donor":      req["donor"],
			"amount":     req["amount"],
			"eventId":    req["eventId"],
		}

		// Save the asset to the ledger
		_, err := assets.PutRecursive(stub, newDonation)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create donation")
		}

		// Return the created donation as JSON
		donationJSON, _ := json.Marshal(newDonation)
		return donationJSON, nil
	},
}
