go run cmd/main.go genesis l2 \
    --deploy-config ../packages/contracts-bedrock/deploy-config/optistable.json \
    --deployment-dir ../packages/contracts-bedrock/deployments/optistable/ \
    --outfile.l2 genesis.json \
    --outfile.rollup rollup.json \
    --l1-rpc $ETH_RPC_URL
