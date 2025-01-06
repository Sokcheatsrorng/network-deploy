package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/msp"
    "github.com/hyperledger-labs/cc-tools/errors"
	// "github.com/golang/protobuf/proto"
	"fmt"
)

// ToJSON converts an object to JSON string
// func ToJSON(data interface{}) ([]byte, error) {
// 	return json.Marshal(data)
// }

// create toJson Return as ICCError
func ToJSON(data interface{}) ([]byte, errors.ICCError){
    assetJSON, err := json.Marshal(data)

    if err != nil {
        return nil,errors.WrapError(err,"Failed t omarshal dasta to JSON")
    }
    return assetJSON,nil
}

// FromJSON converts JSON  to  object
func FromJSON(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}


func GetCreatorMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
    
    creator, err := ctx.GetStub().GetCreator()
    
    if err != nil {
        return "", fmt.Errorf("error getting creator: %v", err)
    }

    sID := &msp.SerializedIdentity{}

    if err := proto.Unmarshal(creator, sID); err != nil {
        return "", fmt.Errorf("error unmarshaling SerializedIdentity: %v", err)
    }

    return sID.Mspid, nil
}

func IsValidRecipientForMSP(recipient, mspID string) bool {
    // Add logic to check recipient's organization
    // logic to validate the recipient and mspID 



    return true // Replace with actual validation
}

// func ToJSON(data interface{}) ([]byte, error) {
//     return json.Marshal(data)
// }
// generate the The uniques id 
func GenereteHashID(donor,donationEventID string) string{
    data := donor + donationEventID 
    hash := sha512.Sum384([]byte(data))
    return hex.EncodeToString(hash[:])
}


func ValidateOrg(ctx contractapi.TransactionContextInterface, allowedOrgs ...string) error {
    // Retrieve the serialized identity of the client
    creator, err := ctx.GetStub().GetCreator()
    if err != nil {
        return fmt.Errorf("failed to get client identity: %v", err)
    }

    // Deserialize the serialized identity
    serializedIdentity := &msp.SerializedIdentity{}
    if err = proto.Unmarshal(creator, serializedIdentity); err != nil {
        return fmt.Errorf("failed to deserialize identity: %v", err)
    }

    // extract MSP ID from the deserialized identity
    mspID := serializedIdentity.Mspid

    // check if MSP ID is in the list of allowed or Not  organizations
    for _, org := range allowedOrgs {
        if mspID == org {
            return nil // Organization is valid 
        }
    }

    // Return an error if the MSP ID is not in the allowed list
    return fmt.Errorf("access denied: %s is not in the allowed organizations: %v", mspID, allowedOrgs)
}




