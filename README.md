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

## Run multiple Nodes locally
1. Create a initial Blockchain:

    `coin init --address <address for genesis reward>`
2. Copy the Blockchain for each node:

     `cp blockchain_1.db blockchain_3000.db`
3. Start the central Node. Central node is running at `localhost:3000`:

    `coin server --address <address for miner rewards> --node 3000`
4. Start miner node:

    `coin server --address <address for miner rewards> --node 3001`

5. To access the wallet of a specific node add the parameter `--node <node id>`:

    `coin list --node 3001`