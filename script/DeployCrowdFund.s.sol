// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script} from "forge-std/Script.sol";
import {CrowdFund} from "../src/CrowdFund.sol";

contract DeployCrowdFundScript is Script {
    function run() external returns (CrowdFund crowdFund) {
        vm.startBroadcast();
        crowdFund = new CrowdFund();
        vm.stopBroadcast();
    }
}
