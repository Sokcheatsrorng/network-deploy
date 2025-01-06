package txdefs

import (
	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// ReadDonationEvent retrieves a donation event by its ID.
var ReadDonationEvent = tx.Transaction{
	Tag:         "readDonationEvent",
	Label:       "Read Donation Event",
	Description: "Retrieve a donation event by its ID",
	Method:      "GET",
	Callers: []accesscontrol.Caller{
		//  can read by Org1 and Org2 
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
			Label:       "Event ID",
			Description: "Unique identifier of the donation event",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// Extract arguments
		id, _ := req["id"].(string)

		// Retrieve asset from the ledger
		eventAsset, err := stub.GetState(id)
		if err != nil || eventAsset == nil {
			return nil, errors.NewCCError("Donation event not found", 404)
		}

		// Return the event as JSON
		return eventAsset, nil
	},
}
