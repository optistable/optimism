## Prerequisites
You’ll need the following software installed to follow this tutorial:

- Git
- Go
- Node
- Pnpm
- Foundry
- Make
- jq
- direnv

### Build the Optimism Monorepo

1. Run `pnpm install`.
2. Run
```
make op-node op-batcher op-proposer
pnpm build
```

### Generation of l2 files (genesis.json, rollup.json and jwt.txt)

Please consider using this command for sepolia instead of the provided on the getting started guide:
```
go run cmd/main.go genesis l2 \
    --deploy-config ../packages/contracts-bedrock/deploy-config/sepolia.json \
    --deployment-dir ../packages/contracts-bedrock/deployments/sepolia/ \
    --outfile.l2 genesis.json \
    --outfile.rollup rollup.json \
    --l1-rpc <RPC>
```


## How to start the rollup after being deployed

#### 1. Reinitialize op-geth

- Create datadir directory if it was no created before `mkdir datadir`
- Delete the geth data. `rm -rf datadir/geth`
- Rerun init. `build/bin/geth init --datadir=datadir genesis.json`
- Run op-geth with the following commands:
```
./build/bin/geth \
        --datadir ./datadir \
        --http \
        --http.corsdomain="*" \
        --http.vhosts="*" \
        --http.addr=0.0.0.0 \
        --http.api=web3,debug,eth,txpool,net,engine \
        --ws \
        --ws.addr=0.0.0.0 \
        --ws.port=8546 \
        --ws.origins="*" \
        --ws.api=debug,eth,txpool,net,engine \
        --syncmode=full \
        --gcmode=archive \
        --nodiscover \
        --maxpeers=0 \
        --networkid=42069 \
        --authrpc.vhosts="*" \
        --authrpc.addr=0.0.0.0 \
        --authrpc.port=8551 \
        --authrpc.jwtsecret=./jwt.txt \
        --rollup.disabletxpoolgossip=true
```

#### 2. Run op-node

```
./bin/op-node \
	--l2=http://localhost:8551 \
	--l2.jwt-secret=./jwt.txt \
	--sequencer.enabled \
	--sequencer.l1-confs=3 \
	--verifier.l1-confs=3 \
	--rollup.config=./rollup.json \
	--rpc.addr=0.0.0.0 \
	--rpc.port=8547 \
	--p2p.disable \
	--rpc.enable-admin \
	--p2p.sequencer.key=$SEQ_KEY \
	--l1=https://eth-sepolia.g.alchemy.com/v2/WQznNJD41WwbDELqAUoVyVDaQT5G-Q79 \
	--l1.rpckind=alchemy
```

#### 3. Run op-batcher

```
./bin/op-batcher \
    --l2-eth-rpc=http://localhost:8545 \
    --rollup-rpc=http://localhost:8547 \
    --poll-interval=1s \
    --sub-safety-margin=6 \
    --num-confirmations=1 \
    --safe-abort-nonce-too-low-count=3 \
    --resubmission-timeout=30s \
    --rpc.addr=0.0.0.0 \
    --rpc.port=8548 \
    --rpc.enable-admin \
    --max-channel-duration=1 \
    --l1-eth-rpc=$L1_RPC \
    --private-key=$BATCHER_KEY
```
#### 4. Run op-proposer

```
./bin/op-proposer \
    --poll-interval=12s \
    --rpc.port=8560 \
    --rollup-rpc=http://localhost:8547 \
    --l2oo-address=$L2OO_ADDR \
    --private-key=$PROPOSER_KEY \
    --l1-eth-rpc=$L1_RPC
```


# Notes on the getting started guide

Regarding the guide https://stack.optimism.io/docs/build/getting-started/ please look at the following notes:

- deploy-config/getting-started.json should also have the following field `{ "l1BlockTime": 2 }` Blocktime of l1 should
always be higher or equal than the one of l2. Also if sepolia set `{ "l1ChainID": 11155111 }`
- Please use `direnv allow .` to load .envrc variables
- The SEQ_KEY should be the private key but without 0x at the beginning