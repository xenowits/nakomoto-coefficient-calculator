## Nakamoto Coefficient Calculator

### AIM

The aim of this project is to calculate the nakamoto coefficients for various popular blockchains.

Nakamoto coefficient is a way to calculate the level of decentralization in a particular chain.

### Programming Languages

Golang

### Steps to run
1. Build docker image
```shell
docker build . --platform=linux/amd64 -t xenowits/nc-calc:v0.1.0
```
2. Run the image
```shell
docker run --rm -v "$(pwd):/opt/xenowits" -p 8080:8080 xenowits/nc-calc:v0.1.0
```

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
16. [Nano](https://nano.org/)

### Notes

The actual logic is present inside `/core`. So, ideally a cron job would be run after every `JOB_INTERVAL` which would save/refresh the nakamoto-coefficients `database`.

### Future Work

To add support for multiple other chains in `/v1`
