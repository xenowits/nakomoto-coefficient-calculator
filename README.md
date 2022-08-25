## Nakamoto Coefficient Calculator

### AIM

The aim of this project is to calculate the nakamoto coefficients for various popular blockchains.

Nakamoto coefficient is a way to calculate the level of decentralization in a particular chain.

### Programming Languages

Golang

### Steps to run

1. Install [golang](https://go.dev/doc/install)
2. Also, install [postgresql](https://www.postgresql.org/download/) and make sure it is running in the background
3. After postgres is installed, copy the commands in `db/postgres_script.sql` into the terminal
4. Add following to your `~/.bashrc` or `~/.zshrc`:
   1. `export SOLANA_API_KEY=api_key`
   2. `export DATABASE_URL=postgres://username:password@localhost:5432/postgres`
5. In a separate terminal, run `go run core/main.go`. This will start the core logic of calculating the nakamoto coefficients.
6. If you want to start the server, run `go run server/main.go` in another terminal.

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

### Notes

1. The actual logic is present inside `/core`. So, ideally a cron job would be run after every `JOB_INTERVAL` which would save/refresh the nakamoto-coefficients `database`.
2. The server code resides inside `/server`. It is a simple server which would only respond to `GET /nakamoto-coefficients`. It basically queries the database and returns the values.

### Future Work

To add support for multiple other chains in `/v1`
