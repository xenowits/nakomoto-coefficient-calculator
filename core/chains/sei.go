package chains

func Sei() (int, error) {
	validatorsURL := "https://rest.sei-apis.com/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=100&status=BOND_STATUS_BONDED"
	poolURL := "https://rest.sei-apis.com/cosmos/staking/v1beta1/pool"
	return FetchCosmosSDKNakaCoeff("sei", validatorsURL, poolURL)
}
