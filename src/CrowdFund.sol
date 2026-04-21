// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract CrowdFund {
    error InvalidGoal();
    error InvalidContribution();
    error InvalidDuration();
    error CampaignNotFound();
    error CampaignEnded();
    error CampaignStillActive();
    error GoalNotReached();
    error GoalAlreadyReached();
    error NotCampaignCreator();
    error NothingToRefund();
    error FundsAlreadyWithdrawn();

    struct Campaign {
        uint256 id;
        address creator;
        string title;
        uint256 goal;
        uint256 pledged;
        uint256 deadline;
        bool withdrawn;
    }

    uint256 public constant MIN_DURATION = 1 hours;
    uint256 public constant MAX_DURATION = 30 days;

    uint256 public nextCampaignId;

    mapping(uint256 => Campaign) public campaigns;
    mapping(uint256 => mapping(address => uint256)) public contributions;

    event CampaignCreated(
        uint256 indexed campaignId, address indexed creator, string title, uint256 goal, uint256 deadline
    );
    event Funded(uint256 indexed campaignId, address indexed funder, uint256 amount);
    event Withdrawn(uint256 indexed campaignId, address indexed creator, uint256 amount);
    event Refunded(uint256 indexed campaignId, address indexed funder, uint256 amount);

    function createCampaign(string calldata title, uint256 goal, uint256 duration)
        external
        returns (uint256 campaignId)
    {
        if (goal == 0) revert InvalidGoal();
        if (duration < MIN_DURATION || duration > MAX_DURATION) revert InvalidDuration();

        campaignId = nextCampaignId++;

        campaigns[campaignId] = Campaign({
            id: campaignId,
            creator: msg.sender,
            title: title,
            goal: goal,
            pledged: 0,
            deadline: block.timestamp + duration,
            withdrawn: false
        });

        emit CampaignCreated(campaignId, msg.sender, title, goal, block.timestamp + duration);
    }

    function fund(uint256 campaignId) external payable {
        Campaign storage campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();
        if (block.timestamp >= campaign.deadline) revert CampaignEnded();
        if (msg.value == 0) revert InvalidContribution();

        campaign.pledged += msg.value;
        contributions[campaignId][msg.sender] += msg.value;

        emit Funded(campaignId, msg.sender, msg.value);
    }

    function withdraw(uint256 campaignId) external {
        Campaign storage campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();
        if (msg.sender != campaign.creator) revert NotCampaignCreator();
        if (block.timestamp < campaign.deadline) revert CampaignStillActive();
        if (campaign.pledged < campaign.goal) revert GoalNotReached();
        if (campaign.withdrawn) revert FundsAlreadyWithdrawn();

        campaign.withdrawn = true;
        uint256 amount = campaign.pledged;

        (bool success,) = payable(campaign.creator).call{value: amount}("");
        require(success, "TRANSFER_FAILED");

        emit Withdrawn(campaignId, campaign.creator, amount);
    }

    function refund(uint256 campaignId) external {
        Campaign storage campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();
        if (block.timestamp < campaign.deadline) revert CampaignStillActive();
        if (campaign.pledged >= campaign.goal) revert GoalAlreadyReached();

        uint256 amount = contributions[campaignId][msg.sender];
        if (amount == 0) revert NothingToRefund();

        contributions[campaignId][msg.sender] = 0;

        (bool success,) = payable(msg.sender).call{value: amount}("");
        require(success, "TRANSFER_FAILED");

        emit Refunded(campaignId, msg.sender, amount);
    }

    function getCampaign(uint256 campaignId) external view returns (Campaign memory) {
        Campaign memory campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();

        return campaign;
    }
}
