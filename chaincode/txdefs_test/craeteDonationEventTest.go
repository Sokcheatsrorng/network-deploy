package txdefs_test

// import (
// 	"chaincode-donation/txdefs"
// 	"encoding/json"
// 	"testing"
// 	"time"

// 	"github.com/hyperledger-labs/cc-tools/mock"
// 	"github.com/hyperledger-labs/cc-tools/stubwrapper"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCreateDonationEvent(t *testing.T) {
// 	// Mock the Fabric stub
// 	mockStub := mock.NewMockStub("donationChaincode", nil)

// 	// Wrap the mock stub with StubWrapper
// 	stub := &stubwrapper.StubWrapper{Stub: mockStub}

// 	// Start a mock transaction
// 	mockStub.MockTransactionStart("testTxID")
// 	defer mockStub.MockTransactionEnd("testTxID")

// 	// Input arguments
// 	args := map[string]interface{}{
// 		"id":          "event123",
// 		"eventName":   "Annual Fundraiser",
// 		"recipient":   "org1MSP",
// 		"description": "An event to raise funds for education.",
// 		"timestamp":   time.Now().Format(time.RFC3339),
// 	}

// 	// Call the CreateDonationEvent routine
// 	result, err := txdefs.CreateDonationEvent.Routine(stub, args)
// 	assert.NoError(t, err, "Transaction should not return an error")
// 	assert.NotNil(t, result, "Result should not be nil")

// 	// Verify that the asset was stored in the state
// 	stateBytes, err := mockStub.GetState("event123")
// 	assert.NoError(t, err, "State retrieval should not return an error")
// 	assert.NotNil(t, stateBytes, "Asset should exist in state")

// 	// Verify the asset's data
// 	var storedAsset map[string]interface{}
// 	err = json.Unmarshal(stateBytes, &storedAsset)
// 	assert.NoError(t, err, "State data should unmarshal without error")
// 	assert.Equal(t, "Annual Fundraiser", storedAsset["eventName"], "Event name should match")
// 	assert.Equal(t, "org1MSP", storedAsset["recipient"], "Recipient should match")
// 	assert.Equal(t, "An event to raise funds for education.", storedAsset["description"], "Description should match")
// }
