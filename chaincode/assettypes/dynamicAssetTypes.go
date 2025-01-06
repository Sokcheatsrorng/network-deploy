package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

// DynamicAssetTypes contains the configuration for the Dynamic AssetTypes feature.

// only the Org1 can create the assetType
var DynamicAssetTypes = assets.DynamicAssetType{
	Enabled:     true,
	AssetAdmins: []string{`org1MSP`, "orgMSP"},
}