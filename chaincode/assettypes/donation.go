package assettypes

import (
	"fmt"

	"github.com/hyperledger-labs/cc-tools/assets"
)


var Donation = assets.AssetType{
    Tag:         "donation",
    Label:       "Donation",
    Description: "Represents a donation",

    Props: []assets.AssetProp{
        {
            IsKey:    true,
            Tag:      "id",
            Label:    "Donation ID",
            DataType: "string",
            Writers:  []string{"org1MSP"},
        },
		{
			Required: true,
			Tag:      "donationEventID",
			Label:    "Donation Event ID",
			DataType: "->donationEvent", // Reference to DonationEvent
		},
        {
            Required: true,
            Tag:      "donor",
            Label:    "Donor",
            DataType: "string",
            // Validate funcion not none or ""
			Validate: func(donor interface{}) error {
				// donorStr := donor.(string)
                if donorStr, ok := donor.(string); !ok || donorStr == "" {
					return fmt.Errorf("Donor must be non-empty")
				}
				return nil
			},
        },
        {
            Required: true,
            Tag:      "amount",
            Label:    "Amount",
            DataType: "number",
            // Validate funcion not none or ""
			Validate: func(amount interface{}) error {
				amountFloat,ok := amount.(float64)
				if !ok || amountFloat <= 0  {
					return fmt.Errorf("Amount must positive number!")
				}
				return nil
			},
        },
        {
            Tag:      "message",
            Label:    "Message",
            DataType: "string",
        },
        {
            Required: true,
            Tag:      "recipient",
            Label:    "Recipient",
            DataType: "string",
        },
        {
            Required: true,
            Tag:      "timestamp",
            Label:    "Timestamp",
            DataType: "datetime",
        },
    },
}
