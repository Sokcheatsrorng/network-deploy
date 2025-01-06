package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

// This function shall make dynamic when donation need to have secrect between Org with Org
// Secret is and information available only for org2 and org3
// Collections.json configuration is necessary
var Secret = assets.AssetType{
	Tag:         "secret",
	Label:       "Secret",
	Description: "Secret between Org2 and Org3",

	Readers: []string{"org2MSP", "org3MSP", "orgMSP"},

	Props: []assets.AssetProp{
		{
			// Primary Key
			IsKey:    true,
			Tag:      "secretName",
			Label:    "Secret Name",
			DataType: "string",
			// This means only org2 can create the asset (org3 can only read)
			Writers:  []string{`org2MSP`, "orgMSP"}, 
		},
		{
			// Mandatory Property
			Required: true,
			Tag:      "secret",
			Label:    "Secret",
			DataType: "string",
		},
	},
}
