package chains

func Osmosis() (int, error) {
	validatorURL := "https://rest.osmosis.goldenratiostaking.net/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	poolURL := "https://rest.osmosis.goldenratiostaking.net/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("osmosis", validatorURL, poolURL)
}
