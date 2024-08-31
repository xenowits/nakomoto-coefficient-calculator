package chains

func Stargaze() (int, error) {
	validatorsURL := "https://rest.stargaze-apis.com/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	stakingPoolURL := "https://rest.stargaze-apis.com/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("stargaze", validatorsURL, stakingPoolURL)
}
