package assettypes

import (
	"fmt"
	"regexp"
	"github.com/hyperledger-labs/cc-tools/assets"
)
// testing evn 

// DonationEvent defines the Donation Event asset in the donation system
var DonationEvent = assets.AssetType{
	Tag:         "donationEvent",
	Label:       "Donation Event",
	Description: "- The purpose of the Donation Event is to Educate Poor Young Generation",

	// Properties of the Donation Event
	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "eventID",
			Label:    "Event ID",
			DataType: "string",
			// Restrict writers dynamically to the invoking organization
			// $org cc-tools provide the functionalty (dynamic form user )
			Writers: []string{`$org`},
		},
		{
			Required: true,
			Tag:      "eventName",
			Label:    "Event Name",
			DataType: "string",
			// Validation function to ensure the name is non-empty
			Validate: func(eventName interface{}) error {
				nameStr := eventName.(string)
				if nameStr == "" {
					return fmt.Errorf("Event Name must be non-empty")
				}
				return nil
			},
			Writers: []string{`$org`},
		},
		{
			Required: true,
			Tag:      "recipient",
			Label:    "Recipient Organization",
			DataType: "string",
			// Validation to ensure the organization is valid
			Validate: func(recipient interface{}) error {
				validRecipients := []string{"org1MSP", "org2MSP"}
				recipientStr := recipient.(string)
				for _, valid := range validRecipients {
					if recipientStr == valid {
						return nil
					}
				}
				return fmt.Errorf("Invalid recipient organization: %s", recipientStr)
			},
			Writers: []string{`$org`},
		},
		{
			Required: false,
			Tag:      "description",
			Label:    "Event Description",
			DataType: "string",
			Writers: []string{`$org`},
		},
		{
			Required: true,
			Tag:      "email",
			Label:    "Email",
			DataType: "string",
			// Enhanced email validation using regex
			Validate: func(email interface{}) error {
				emailStr := email.(string)
				regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
				if matched := regexp.MustCompile(regex).MatchString(emailStr); !matched {
					return fmt.Errorf("Invalid email format")
				}
				return nil
			},
			Writers: []string{`$org`},
		},
		{
			Required: true,
			Tag:      "timestamp",
			Label:    "Timestamp",
			DataType: "datetime",
			Writers: []string{`$org`},
		},
		{
			Required: true,
			Tag:      "organization",
			Label:    "Organization",
			DataType: "string",
			// Validation to check the organization is valid
			Validate: func(org interface{}) error {
				validOrgs := []string{"org1MSP", "org2MSP"}
				orgStr := org.(string)
				for _, validOrg := range validOrgs {
					if orgStr == validOrg {
						return nil
					}
				}
				return fmt.Errorf("Invalid organization: %s", orgStr)
			},
			Writers: []string{`$org`},
		},
		{
			Required: false,
			Tag:      "donations",
			Label:    "Donations",
			DataType: "[]->donation", // Reference to the Donation asset type
			Writers: []string{`$org`},
		},
	},
}
