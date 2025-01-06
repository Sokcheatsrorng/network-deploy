package txdefs

import (
	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// DeleteDonation deletes a donation by its ID.
var DeleteDonation = tx.Transaction{
	Tag:         "deleteDonation",
	Label:       "Delete Donation",
	Description: "Delete a donation record from the ledger",
	Method:      "DELETE",
	Callers: []accesscontrol.Caller{
		// org that's can detete this donation event org1 and org2 
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
			Label:       "Donation ID",
			Description: "Unique identifier for the donation",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// Extract ID
		id, ok := req["id"].(string)

		if !ok {
			return nil,errors.NewCCError("ID is requried and ID must be a string",400)
		}
		// create a key instance of the donation 
		assetKey, err := assets.NewKey(map[string]interface{}{
			"@assetType":"donation",
			"@key":id,
		})

		if err != nil {
			return nil, errors.WrapError(err,"Failed to create donation key")
		}

		// go to get the donation with the key that provide 
		asset,err := assetKey.Get(stub)

		if err != nil || asset == nil{
			return nil,errors.NewCCError("Donation has been not found ",404)
		}
		// delete donation using Delete method 
		response, delErr := asset.Delete(stub)
		if delErr !=nil{
			return nil,errors.WrapError(delErr,"Failed to Delete the Donation")
		}
		// Return success message
		return response, nil
	},
}
