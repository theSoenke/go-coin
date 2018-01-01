# go-coin

Experimental Blockchain in Go

## Quick Start
### Install

    go get github.com/thesoenke/go-coin

### Create a Wallet

    coin wallet

### Initialize Blockchain

    coin init --address <address for genesis reward>

### Send coins

    coin send --from <sender address> --to <receiver address> --amount <coins>

## Run multiple nodes locally
### Create an initial Blockchain

    coin init --address <address for genesis reward>

### Copy the Blockchain for each node

    cp blockchain_1.db blockchain_3000.db
    cp blockchain_1.db blockchain_3001.db

### Start a central node

    coin server --address <address for miner rewards> --node 3000

Central node is running at `localhost:3000`

### Start miner node

    coin server --address <address for miner rewards> --node 3001