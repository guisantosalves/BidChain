// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/AuctionFactory.sol";

// Tudo entre vm.startBroadcast() e
// vm.stopBroadcast() é transmitido pra rede real (ou testnet)
contract AuctionScript is Script {
    function run() external {
        uint256 deployKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployKey);

        new AuctionFactory();

        vm.stopBroadcast();
    }
}
