package chains

func Sei() (int, error) {
	url := "https://rest.sei-apis.com/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	return fetchCosmosSDKNakaCoeff("sei", url)
}
