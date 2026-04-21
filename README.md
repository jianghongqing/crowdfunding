# Foundry Crowdfunding Demo

这是一个适合学习 Foundry 和 Solidity 的众筹项目最小实现。

目标不是一开始就做成“生产级平台”，而是先通过一个完整闭环把下面这些核心概念跑通：

- 合约状态设计
- 资金流转逻辑
- 时间条件控制
- Foundry 测试驱动开发
- 本地部署与脚本执行

## 1. 项目结构

```text
src/
  CrowdFund.sol          核心众筹合约
test/
  CrowdFund.t.sol        Foundry 测试
script/
  DeployCrowdFund.s.sol  部署脚本
```

## 2. 需求拆解

先实现一个学习版 MVP，只保留最重要的 4 个能力：

1. 发起人创建众筹活动
2. 用户在截止时间前捐款
3. 达标后发起人提取资金
4. 未达标时捐款人退款

这版故意不做太多高级功能，比如：

- 平台手续费
- 活动取消
- ERC20 代币众筹
- NFT 回报
- 多阶段里程碑拨款

原因很简单：学习阶段最重要的是先把 ETH 众筹的状态机和资金安全流程理解透。

## 3. 合约设计思路

### 3.1 为什么用 `Campaign` 结构体

每个众筹活动都需要一份独立状态，包括：

- 谁创建的
- 目标金额是多少
- 当前已经筹到多少
- 什么时候截止
- 是否已经被提取

这些字段天然属于同一个业务对象，所以适合放进 `Campaign` 结构体里。

### 3.2 为什么单独存 `contributions`

每个活动下，每个用户的捐款额都要能被查询和退款。

所以使用：

```solidity
mapping(uint256 => mapping(address => uint256)) public contributions;
```

第一层 key 是活动 ID，第二层 key 是用户地址，value 是该用户给这个活动捐了多少 ETH。

### 3.3 为什么不用“管理员手动结算”

我们直接把规则写进合约：

- 到期且达标 -> `withdraw()`
- 到期且未达标 -> `refund()`

这样更贴近智能合约“代码即规则”的思想，也更容易测试。

## 4. 开发流程

推荐你以后继续按这个顺序练习：

1. 先写需求和状态图
2. 再定义存储结构和事件
3. 先写 happy path 测试
4. 再补各种 revert 场景
5. 最后写部署脚本和交互命令

这个仓库已经按这种顺序搭好了第一版。

## 5. 常用命令

当前机器上的 Foundry 安装路径在 Git Bash 环境下可直接使用：

```bash
~/.foundry/bin/forge build
~/.foundry/bin/forge test -vv
~/.foundry/bin/forge fmt
~/.foundry/bin/anvil
```

### 本地启动链

```bash
~/.foundry/bin/anvil
```

### 编译

```bash
~/.foundry/bin/forge build
```

### 测试

```bash
~/.foundry/bin/forge test -vv
```

### 部署到本地 Anvil

先启动 `anvil`，再开另一个终端执行：

```bash
~/.foundry/bin/forge script script/DeployCrowdFund.s.sol:DeployCrowdFundScript \
  --rpc-url http://127.0.0.1:8545 \
  --private-key <ANVIL_PRIVATE_KEY> \
  --broadcast
```

## 6. 第一版代码你应该重点看什么

### `createCampaign`

关注两件事：

- 参数校验：目标金额不能为 0，持续时间不能太短或太长
- 状态初始化：为新活动分配 ID，并写入链上存储

### `fund`

关注三件事：

- 活动必须存在
- 还没到截止时间
- 用户转进来的 `msg.value` 会同时更新总筹款额和个人捐款记录

### `withdraw`

这是“成功路径”：

- 只有发起人能提取
- 必须已经到期
- 必须达到目标金额
- 提取前先更新状态，再转账，避免重复提取风险

### `refund`

这是“失败路径”：

- 必须已经到期
- 必须没有达标
- 只能退回自己捐过的金额
- 退款前先把用户记录清零，再转账

## 7. 下一步适合继续做什么

如果你想继续把这个项目当成学习主线，我建议下一阶段按下面顺序扩展：

1. 增加 `cancelCampaign`
2. 增加活动描述、图片 CID 等元数据字段
3. 支持 ERC20 众筹
4. 增加平台手续费和项目方收益分配
5. 接一个简单前端页面

## 8. 参考资料

- [Foundry Book](https://book.getfoundry.sh/)
- [Foundry Installation](https://www.getfoundry.sh/introduction/installation)
