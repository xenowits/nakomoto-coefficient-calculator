## Nakamoto Coefficient Calculator

### AIM

The aim of this project is to calculate the nakamoto coefficients for various popular blockchains.

[Nakamoto coefficient](https://news.earn.com/quantifying-decentralization-e39db233c28e) is a way to calculate the level of decentralization in a particular chain.

Read this amazing [messari report](https://messari.io/report/evaluating-validator-decentralization-geographic-and-infrastructure-distribution-in-proof-of-stake-networks) on operational decentralization of Proof-of-stake networks.

#### Disclaimer

Please note that the values should be interpreted with context since the same objective treatment is applied for all the chains included here, ie,
we simply calculate:
```markdown
nakamoto-coefficient: no of validators controlling 33% of the total network stake
```

Note that the threshold may be different for some blockchains, for example, 50%.
So, I would suggest users to understand the context, cross-verify and examine the results. For any feedback, please join this [discord](https://discord.gg/Una8qmFg).

### Programming Languages

Golang

### Steps to run
1. Build docker image
```shell
docker build . --platform=linux/amd64 -t xenowits/nc-calc:v0.1.4
```
2. Run the image
```shell
docker run --rm \
-e "SOLANA_API_KEY=<YOUR_SOLANA_API_KEY_HERE>" \
-e "RATED_API_KEY=<YOUR_RATED_API_KEY_HERE>" \
-p 8080:8080 xenowits/nc-calc:v0.1.4
```

NOTE: You can get your API Key by signing up [here](https://www.validators.app/users/sign_up?locale=en&network=mainnet).

### Chains currently supported

1. [Agoric](https://agoric.com/)
2. [Aptos](https://aptosfoundation.org/)
3. [Avalanche](https://www.avax.network/)
4. [BNB Smart Chain](https://www.bnbchain.org)
5. [Cardano](https://cardano.org/)
6. [Celestia](https://celestia.org/)
7. [Cosmos](https://cosmos.network/)
8. [Graph Protocol](https://thegraph.com/)
9. [Hedera](https://hedera.com/)
10. [Juno](https://www.junonetwork.io/)
11. [Mina](https://minaprotocol.com/)
12. [MultiversX](https://multiversx.com/)
13. [Near](https://near.org/)
14. [Osmosis Zone](https://osmosis.zone/)
15. [Polygon](https://polygon.technology/)
16. [Polkadot](https://polkadot.network/)
17. [Pulsechain](https://pulsechain.com/)
18. [Regen Network](https://www.regen.network/)
19. [Sei](https://sei.io/)
20. [Solana](https://solana.com/)
21. [Stargaze](https://stargaze.zone/)
22. [Sui](https://sui.io/)
23. [Terra](https://www.terra.money/)
24. [Thorchain](https://www.thorchain.com/)
25. [Nano](https://www.nano.org)

### Notes

The actual logic is present inside `/core`. A goroutine runs every 6 hours which updates the nakamoto coefficients for all the chains.

### Future Work

To add support for multiple other chains in `/v1`.