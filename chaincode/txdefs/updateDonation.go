package txdefs

import (
	"chaincode-donation/utils"
	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	// "github.com/hyperledger/fabric-protos-go/msp"
)

// UpdateDonation updates an existing donation record.
var UpdateDonation = tx.Transaction{
	Tag:         "updateDonation",
	Label:       "Update Donation",
	Description: "Update an existing donation record",
	Method:      "PUT",
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
			Label:       "Donation ID",
			Description: "Unique identifier for the donation",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "donor",
			Label:       "Donor",
			Description: "Updated donor name or organization",
			DataType:    "string",
		},
		{
			Tag:         "amount",
			Label:       "Amount",
			Description: "Updated donation amount",
			DataType:    "float",
		},
		{
			Tag:         "eventId",
			Label:       "Event ID",
			Description: "Updated related event ID",
			DataType:    "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// Extract ID
		id, ok := req["id"].(string)
		//  checking the id 
		if !ok || id == "" {
			return nil,errors.NewCCError("Invaild ID, Please the id again ",400) 
			
		}

		// create a assset key to retrieve the doantion 
		assetKey, err := assets.NewKey(map[string]interface{}{
			"@assetType":"donation",
			"@key": id,
		})

		// handle when fail to create key 
		if err != nil{
			return nil, errors.WrapError(err,"Failer to create asset key")
		}
		// retrieve or get the exsiting doantion from couchDB 
		asset, err := assetKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err,"Failed to retrive Doantion")
		}
		if asset == nil {
			return nil, errors.NewCCError("Donation not found", 404)
		}

		// Update fields donation 
		// var asset *assets.Asset

		// Update fields in the asset

		assetMap := (map[string]interface{})(*asset)

		if donor, ok := req["donor"].(string); ok {
			assetMap["donor"] = donor
		}
		if amount, ok := req["amount"].(float64); ok {
			assetMap["amount"] = amount
		}
		if eventId, ok := req["eventId"].(string); ok {
			assetMap["eventId"] = eventId
		}

		// Save updated asset
		udateAsset, err := asset.Update(stub,assetMap)

		if err != nil {
			return nil, errors.WrapError(err, "Failed to update donation")
		}

		// Return the updated asset as JSON
		updatedAssetJSON, err := utils.ToJSON(udateAsset)
		if err != nil {
			return nil, errors.WrapError(err,"Failed to update the assets to JSON ")
		}
		// return the asset that have been update 
		return updatedAssetJSON, nil
	},
}
