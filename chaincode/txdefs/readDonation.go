package txdefs

import (
	"chaincode-donation/utils"
	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// ReadDonation retrieves a donation by its ID.
var ReadDonation = tx.Transaction{
	Tag:         "readDonation",
	Label:       "Read Donation",
	Description: "Retrieve a donation by its ID",
	Method:      "GET",
	Callers: []accesscontrol.Caller{
		{
			// can get read this only org1 and org2 
			MSP: "OrgMSP1",
		},
		{
			MSP: "OrgMSP2",
		},
		{
			MSP: "OrgMSP",
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
		// Retrieve asset
		id, ok  := req["id"].(string)
		if !ok || id ==""{
			return nil,errors.NewCCError("Invailed or missing the 'id' argument",400)

		}
		// carete a instance key to point  donation 
		assetKey, err := assets.NewKey(map[string]interface{}{
			"@asssetType":"donation",
			"@key": id,
		})

		if err != nil {
			return nil, errors.WrapError(err,"Failed to create a asset key")
		}

		// Retrieve the donaton usinge the created key 
		asset, err := assetKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err,"Failed to retrieve donation")
		}

		if asset == nil {
			return nil, errors.NewCCError("Doantion  has been not found",404)
		}
		// Marshal the asset to JSON and return it
		assetJSON, err := utils.ToJSON(asset)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to marshal asset to JSON")
		}

		// Return the marshaled JSON as the response
		return assetJSON, nil
	},
}
