## Util commands

- Run  ``cat deployments/getting-started/L1StandardBridgeProxy.json | jq -r .address`` under packages/contracts-bedrock/ to get the address of the rollup.
- `cast call --rpc-url http://localhost:8545  contractAddress "function()"` to call a read function of a deployed contract. Add `` | cast to-ascii``to parse the result.
- ``cast send --rpc-url http://localhost:8545  --private-key $PRIVATE_KEY contractAddress "function(uint64,uint64)" arg1 arg2`` to call a write function of a deployed contract
- ``cast send --private-key $PRIVATE_KEY --rpc-url http://localhost:8545 recipientAddress  --value wei`` to send eth to another address
- ``cast balance --rpc-url http://localhost:8545 address`` to check the goerli balance of an address
- ``cast block-number --rpc-url http://localhost:8545 ``to check current block number

# Notes on the getting started guide

Regarding the guide https://stack.optimism.io/docs/build/getting-started/ please look at the following notes:

- deploy-config/getting-started.json should also have the following field `{ "l1BlockTime": 2 }` Blocktime of l1 should
always be higher or equal than the one of l2. Also if sepolia set `{ "l1ChainID": 11155111 }`
- Please use `direnv allow .` to load .envrc variables
- The SEQ_KEY should be the private key but without 0x at the beginning
- For op-geth be careful when deleting the datadir directoy, l2 blocks could come back to the first blocks.