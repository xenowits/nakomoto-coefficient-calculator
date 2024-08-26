package chains

func Regen() (int, error) {
	validatorURL := "https://regen.api.m.stavr.tech/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=100&status=BOND_STATUS_BONDED"
	poolURL := "https://regen.api.m.stavr.tech/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("regen", validatorURL, poolURL)
}
