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
3. [Avail](https://www.availproject.org/)
4. [Avalanche](https://www.avax.network/)
5. [BNB Smart Chain](https://www.bnbchain.org)
6. [Cardano](https://cardano.org/)
7. [Celestia](https://celestia.org/)
8. [Cosmos](https://cosmos.network/)
9. [Graph Protocol](https://thegraph.com/)
10. [Hedera](https://hedera.com/)
11. [Juno](https://www.junonetwork.io/)
12. [Mina](https://minaprotocol.com/)
13. [MultiversX](https://multiversx.com/)
14. [Nano](https://nano.org/)
15. [Near](https://near.org/)
16. [Osmosis Zone](https://osmosis.zone/)
17. [Polygon](https://polygon.technology/)
18. [Polkadot](https://polkadot.network/)
19. [Pulsechain](https://pulsechain.com/)
20. [Regen Network](https://www.regen.network/)
21. [Sei](https://sei.io/)
22. [Solana](https://solana.com/)
23. [Stargaze](https://stargaze.zone/)
24. [Sui](https://sui.io/)
25. [Terra](https://www.terra.money/)
26. [Tezos](https://tezos.com/)
27. [Thorchain](https://www.thorchain.com/)

### Notes

The actual logic is present inside `/core`. A goroutine runs every 6 hours which updates the nakamoto coefficients for all the chains.

### Future Work

To add support for multiple other chains in `/v1`.