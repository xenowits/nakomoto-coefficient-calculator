package chains

func Agoric() (int, error) {
	validatorURL := "https://main.api.agoric.net/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=100&status=BOND_STATUS_BONDED"
	stakingPoolURL := "https://main.api.agoric.net/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("agoric", validatorURL, stakingPoolURL)
}
