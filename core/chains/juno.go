package chains

func Juno() (int, error) {
	validatorsURL := "https://api.juno.basementnodes.ca/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=100&status=BOND_STATUS_BONDED"
	stakingPoolURL := "https://api.juno.basementnodes.ca/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("juno", validatorsURL, stakingPoolURL)
}
