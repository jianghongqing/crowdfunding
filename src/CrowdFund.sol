// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract CrowdFund {
    // --------
    // Errors
    // --------
    // Custom errors 比字符串 revert 更省 gas，也更容易在测试里精准断言。
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

    /// @notice 单个众筹活动的核心状态。
    struct Campaign {
        uint256 id;
        address creator;
        string title;
        uint256 goal;
        uint256 pledged;
        uint256 deadline;
        bool withdrawn;
    }

    /// @notice 活动持续时间限制，避免无限期活动或极短活动。
    uint256 public constant MIN_DURATION = 1 hours;
    uint256 public constant MAX_DURATION = 30 days;

    /// @notice 下一个可分配的活动 ID（递增）。
    uint256 public nextCampaignId;

    /// @notice campaignId => Campaign
    mapping(uint256 => Campaign) public campaigns;
    /// @notice campaignId => user => amount，记录每个用户在某个活动的捐款额。
    mapping(uint256 => mapping(address => uint256)) public contributions;

    event CampaignCreated(
        uint256 indexed campaignId, address indexed creator, string title, uint256 goal, uint256 deadline
    );
    event Funded(uint256 indexed campaignId, address indexed funder, uint256 amount);
    event Withdrawn(uint256 indexed campaignId, address indexed creator, uint256 amount);
    event Refunded(uint256 indexed campaignId, address indexed funder, uint256 amount);

    /// @notice 创建新的众筹活动。
    /// @param title 活动标题
    /// @param goal 目标金额（wei）
    /// @param duration 持续时间（秒）
    /// @return campaignId 新活动 ID
    function createCampaign(string calldata title, uint256 goal, uint256 duration)
        external
        returns (uint256 campaignId)
    {
        if (goal == 0) revert InvalidGoal();
        if (duration < MIN_DURATION || duration > MAX_DURATION) revert InvalidDuration();

        campaignId = nextCampaignId++;

        // 这里把关键字段一次性初始化，形成可验证的链上状态快照。
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

    /// @notice 向活动捐款（ETH）。
    /// @dev 仅在截止前允许捐款。
    function fund(uint256 campaignId) external payable {
        Campaign storage campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();
        if (block.timestamp >= campaign.deadline) revert CampaignEnded();
        if (msg.value == 0) revert InvalidContribution();

        campaign.pledged += msg.value;
        contributions[campaignId][msg.sender] += msg.value;

        emit Funded(campaignId, msg.sender, msg.value);
    }

    /// @notice 达标后由活动创建者提取资金。
    /// @dev 安全顺序：先更新状态，再外部转账，避免重复提取风险。
    function withdraw(uint256 campaignId) external {
        Campaign storage campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();
        if (msg.sender != campaign.creator) revert NotCampaignCreator();
        if (block.timestamp < campaign.deadline) revert CampaignStillActive();
        if (campaign.pledged < campaign.goal) revert GoalNotReached();
        if (campaign.withdrawn) revert FundsAlreadyWithdrawn();

        // Checks-Effects-Interactions: 先写状态，后与外部地址交互。
        campaign.withdrawn = true;
        uint256 amount = campaign.pledged;

        (bool success,) = payable(campaign.creator).call{value: amount}("");
        require(success, "TRANSFER_FAILED");

        emit Withdrawn(campaignId, campaign.creator, amount);
    }

    /// @notice 活动失败后，捐款人按自己贡献金额退款。
    /// @dev 安全顺序：先把可退款额度清零，再转账，防止重复退款。
    function refund(uint256 campaignId) external {
        Campaign storage campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();
        if (block.timestamp < campaign.deadline) revert CampaignStillActive();
        if (campaign.pledged >= campaign.goal) revert GoalAlreadyReached();

        uint256 amount = contributions[campaignId][msg.sender];
        if (amount == 0) revert NothingToRefund();

        // 先清零，防重入/重复调用。
        contributions[campaignId][msg.sender] = 0;

        (bool success,) = payable(msg.sender).call{value: amount}("");
        require(success, "TRANSFER_FAILED");

        emit Refunded(campaignId, msg.sender, amount);
    }

    /// @notice 查询活动完整信息。
    function getCampaign(uint256 campaignId) external view returns (Campaign memory) {
        Campaign memory campaign = campaigns[campaignId];
        if (campaign.creator == address(0)) revert CampaignNotFound();

        return campaign;
    }
}
