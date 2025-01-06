package txdefs

import (
	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// DeleteDonationEvent deletes a donation event by its ID.
var DeleteDonationEvent = tx.Transaction{
	Tag:         "deleteDonationEvent",
	Label:       "Delete Donation Event",
	Description: "Delete a donation event from the ledger",
	Method:      "DELETE",
	Callers: []accesscontrol.Caller{
		{
			MSP: "OrgMSP1",
			OU:  "admin",
		},
		{
			MSP: "OrgMSP2",
			OU:  "admin",
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
		// Extract the ID(checking the id first )
		id, ok := req["id"].(string)
		if !ok || id == "" {
			return nil, errors.NewCCError("Invalid or missing 'id' argument", 400)
		}

		// Create a key instance for the donation event
		assetKey, err := assets.NewKey(map[string]interface{}{
			"@assetType": "donationEvent",
			"@key":       id,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create Donation Event key")
		}

		// Retrieve the asset to confirm its existence
		asset, err := assetKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to retrieve donation event")
		}
		if asset == nil {
			return nil, errors.NewCCError("Donation event not found", 404)
		}

		// Delete the asset (Delete function inside in the cc-tools )
		response, err := asset.Delete(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to delete donation event")
		}

		return response, nil
	},
}
