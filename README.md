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

1. [Cosmos](https://cosmos.network/)
2. [Polygon](https://polygon.technology/)
3. [Binance Smart Chain](https://www.binance.com)
4. [Osmosis Zone](https://osmosis.zone/)
5. [Mina](https://minaprotocol.com/)
6. [Solana](https://solana.com/)
7. [Avalanche](https://www.avax.network/)
8. [Terra](https://www.terra.money/)
9. [Graph Protocol](https://thegraph.com/)
10. [Thorchain](https://www.thorchain.com/)
11. [Near](https://near.org/)
12. [Juno](https://www.junonetwork.io/)
13. [Ethereum 2](https://ethereum.org/)
14. [Regen Network](https://www.regen.network/)
15. [Agoric](https://agoric.com/)
16. [Stargaze](https://stargaze.zone/)
17. [Hedera](https://hedera.com/)
18. [Sui](https://sui.io/)
19. [Pulsechain](https://pulsechain.com/)
20. [Celestia](https://celestia.org/)
21. [MultiversX](https://multiversx.com/)

### Notes

The actual logic is present inside `/core`. A goroutine runs every 6 hours which updates the nakamoto coefficients for all the chains.

### Future Work

To add support for multiple other chains in `/v1`.
