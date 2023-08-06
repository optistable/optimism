/// SPDX-License-Identifier: MIT
pragma solidity =0.8.15;

import "forge-std/Script.sol";
import "src/L1/L1Burn.sol";


/// To deploy this script run:
/// forge script scripts/DeployL1Burn.s.sol:DeployL1Burn --private-key $PRIVATE_KEY --broadcast --rpc-url http://localhost:8545
contract DeployL1Burn is Script {
    L1Burn public l1Burn;

    function run() public {
        vm.startBroadcast();

        l1Burn = new L1Burn();

        vm.stopBroadcast();
    }
}