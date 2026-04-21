// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {CrowdFund} from "../src/CrowdFund.sol";

contract CrowdFundTest is Test {
    CrowdFund internal crowdFund;

    address internal creator = makeAddr("creator");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");

    uint256 internal constant GOAL = 10 ether;
    uint256 internal constant DURATION = 3 days;

    function setUp() public {
        crowdFund = new CrowdFund();

        vm.deal(creator, 20 ether);
        vm.deal(alice, 20 ether);
        vm.deal(bob, 20 ether);
    }

    function test_CreateCampaign() public {
        vm.prank(creator);
        uint256 campaignId = crowdFund.createCampaign("Build an onchain gadget", GOAL, DURATION);

        CrowdFund.Campaign memory campaign = crowdFund.getCampaign(campaignId);

        assertEq(campaign.id, 0);
        assertEq(campaign.creator, creator);
        assertEq(campaign.title, "Build an onchain gadget");
        assertEq(campaign.goal, GOAL);
        assertEq(campaign.pledged, 0);
        assertEq(campaign.deadline, block.timestamp + DURATION);
        assertFalse(campaign.withdrawn);
    }

    function test_FundCampaign() public {
        uint256 campaignId = _createCampaign();

        vm.prank(alice);
        crowdFund.fund{value: 2 ether}(campaignId);

        CrowdFund.Campaign memory campaign = crowdFund.getCampaign(campaignId);

        assertEq(campaign.pledged, 2 ether);
        assertEq(crowdFund.contributions(campaignId, alice), 2 ether);
    }

    function test_WithdrawAfterGoalReached() public {
        uint256 campaignId = _createCampaign();

        vm.prank(alice);
        crowdFund.fund{value: 6 ether}(campaignId);

        vm.prank(bob);
        crowdFund.fund{value: 4 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        uint256 creatorBalanceBefore = creator.balance;

        vm.prank(creator);
        crowdFund.withdraw(campaignId);

        CrowdFund.Campaign memory campaign = crowdFund.getCampaign(campaignId);

        assertTrue(campaign.withdrawn);
        assertEq(creator.balance, creatorBalanceBefore + GOAL);
    }

    function test_RefundAfterCampaignFails() public {
        uint256 campaignId = _createCampaign();

        vm.prank(alice);
        crowdFund.fund{value: 3 ether}(campaignId);

        vm.prank(bob);
        crowdFund.fund{value: 2 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        uint256 aliceBalanceBefore = alice.balance;

        vm.prank(alice);
        crowdFund.refund(campaignId);

        assertEq(alice.balance, aliceBalanceBefore + 3 ether);
        assertEq(crowdFund.contributions(campaignId, alice), 0);
    }

    function test_RevertIfFundingAfterDeadline() public {
        uint256 campaignId = _createCampaign();

        vm.warp(block.timestamp + DURATION);

        vm.prank(alice);
        vm.expectRevert(CrowdFund.CampaignEnded.selector);
        crowdFund.fund{value: 1 ether}(campaignId);
    }

    function test_RevertIfNonCreatorWithdraws() public {
        uint256 campaignId = _createCampaign();

        vm.prank(alice);
        crowdFund.fund{value: GOAL}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.prank(alice);
        vm.expectRevert(CrowdFund.NotCampaignCreator.selector);
        crowdFund.withdraw(campaignId);
    }

    function test_RevertIfRefundTwice() public {
        uint256 campaignId = _createCampaign();

        vm.prank(alice);
        crowdFund.fund{value: 1 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.startPrank(alice);
        crowdFund.refund(campaignId);
        vm.expectRevert(CrowdFund.NothingToRefund.selector);
        crowdFund.refund(campaignId);
        vm.stopPrank();
    }

    function _createCampaign() internal returns (uint256 campaignId) {
        vm.prank(creator);
        campaignId = crowdFund.createCampaign("Build an onchain gadget", GOAL, DURATION);
    }
}
