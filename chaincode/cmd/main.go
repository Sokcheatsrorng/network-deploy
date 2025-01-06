package main

import (
	"chaincode-donation/contract"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
		// Create an instance of DonationContract
		donationContract := new(contract.DonationContract)
	
		// Create the chaincode
		chaincode, err := contractapi.NewChaincode(donationContract)
		
		if err != nil {
			log.Panicf("Error creating donation chaincode: %v", err)
		}
	
		// Start the chaincode
		if err := chaincode.Start(); err != nil {
			log.Panicf("Error starting donation chaincode: %v", err)
		}
}