# Meteor KV

A simple key-value store built with CometBFT.

This is the simplest possible application for Comet BFT that I could think of.
It only shows the basics ABCI features:

- CheckTx
- FinalizeBlock
- Commit
- Query

## Architecture

This application uses CometBFT.

CometBFT runs the p2p network for connecting to other instances, and a HTTP
server exposing some API that can be used to query the state of the
blockchain or to submit new transactions.

Note: to CometBFT, all the content of the transactions are `[]byte`. It's up
to the application to decide what to do with them.

The non-boilerplate code is in `app.go` and `tx.go`.

The only transaction (tx) accepted is a `Set` transaction, which sets a key
to a certain value. The format is `key=value`. The key must be at least 3
characters long, and must composed by ASCII chars only.

The application uses [BadgerDB v4](https://github.com/dgraph-io/badger) to
store the key-value pairs on disk.

The application exposes a query API to retrieve the value of a key.

## Usage

### Init your node

```bash
go run . init
```

This creates a folder under `~/.meteorkv` containing a genesis file, a
default CometBFT configuration, and a private key for your validator.

The genesis file contains the initial state of the blockchain, and the
current node as its sole validator.

### Start your node

```bash
go run . start
```

This starts CometBFT. By default, a new block (even if it's empty) is created
every second.

You can see the `height` increase in the logs:

```
I[2024-12-04|11:41:33.037] finalizing commit of block                   module=consensus height=2 hash=DF6EFFD2D0C48A37B722541EFF76DECB6063D41B3B8C0ACAA5C14B5ADDA9EC2A root= num_txs=0
I[2024-12-04|11:41:33.045] finalized block                              module=state height=2 num_txs_res=0 num_val_updates=0 block_app_hash=
I[2024-12-04|11:41:33.045] executed block                               module=state height=2 app_hash=
I[2024-12-04|11:41:33.049] committed state                              module=state height=2 block_app_hash=
I[2024-12-04|11:41:33.057] indexed block events                         module=txindex height=2
```

### Submit a transaction

Submit a transaction with the content `cometbft=rocks`:

```bash
curl -s 'localhost:26657/broadcast_tx_commit?tx="cometbft=rocks"'
```

If everything works, the response will be:

```json
{
  "jsonrpc": "2.0",
  "id": -1,
  "result": {
    "check_tx": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gas_wanted": "0",
      "gas_used": "0",
      "events": [],
      "codespace": ""
    },
    "tx_result": {
      "code": 0,
      "data": null,
      "log": "",
      "info": "",
      "gas_wanted": "0",
      "gas_used": "0",
      "events": [],
      "codespace": ""
    },
    "hash": "71276C4844CE72F6C6C868541D10923259F5F8DA5716B230555B36AD309D6FD1",
    "height": "14"
  }
}
```

(see `check_tx.code` and `tx_result.code` in the response, `0` means ok)

### Query a key

Query the value of the key `cometbft`:

```bash
curl -s 'localhost:26657/abci_query?data="cometbft"'
```

The response is:

```json
{
  "jsonrpc": "2.0",
  "id": -1,
  "result": {
    "response": {
      "code": 0,
      "log": "exists",
      "info": "",
      "index": "0",
      "key": "Y29tZXRiZnQ=",
      "value": "cm9ja3M=",
      "proofOps": null,
      "height": "0",
      "codespace": ""
    }
  }
}
```

Since `key` and `value` are only bytes for CometBFT, they are returned as
base64-encoded strings.
