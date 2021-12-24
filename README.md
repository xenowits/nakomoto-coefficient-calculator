## Nakamoto Coefficient Calculator

### AIM

The aim of this project is to calculate the nakamoto coefficients for various popular blockchains.

Nakamoto coefficient is a way to calculate the level of decentralization in a particular chain.

### Programming Languages

Golang

### Steps to run

1. Make sure you have go installed in your system
2. Run `go run cmd/main.go`.

### Chains currently supported

1. [Cosmos](https://cosmos.network/)
2. [Polygon](https://polygon.technology/)
3. [Binance Smart Chain](https://www.binance.com)
4. [Osmosis Zone](https://osmosis.zone/)
5. [Mina](https://minaprotocol.com/)
6. [Solana](https://solana.com/)
7. [Avalanche](https://www.avax.network/)
8. [Terra](https://www.terra.money/)

### Notes
1. Client code (for webpage) is present inside `/web`.
2. The actual logic is present inside `/cmd`. So, ideally a cron job would be run after every `JOB_INTERVAL` which would save/refresh the nakamoto-coefficients `database`.
3. The server code resides inside `/server`. It is a simple server which would only respond to `GET /nakamoto-coefficients`. It basically queries the database and returns the values.

### Future Work

To add support for multiple other chains.
