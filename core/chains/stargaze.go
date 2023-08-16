package chains

func Stargaze() (int, error) {
	url := "https://rest.stargaze-apis.com/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	return fetchCosmosSDKNakamotoCoefficient("stargaze", url)
}
