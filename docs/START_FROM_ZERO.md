# 从零启动

本文档面向第一次在本地跑起整个项目的开发者。

## 依赖准备

至少需要：

- Go 1.25+
- Foundry
- Node.js 18+
- MySQL 8+
- 一个可用的 EVM RPC
- 浏览器钱包

## 1. 合约

在仓库根目录执行：

```bash
forge build
forge test -vv
```

如需本地链：

终端 A：

```bash
anvil
```

终端 B：

```bash
forge script script/DeployCrowdFund.s.sol:DeployCrowdFundScript \
  --rpc-url http://127.0.0.1:8545 \
  --private-key <ANVIL_PRIVATE_KEY> \
  --broadcast
```

记录部署出的：

- 合约地址
- 部署区块

## 2. 初始化 MySQL

创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS crowdfunding CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

执行迁移：

```sql
USE crowdfunding;
SOURCE backend/migrations/001_init.sql;
SOURCE backend/migrations/002_add_campaign_status.sql;
```

## 3. 准备链配置

编辑：

- `backend/config/chain.testnet.example.json`

至少确认这些字段：

- `rpcHttpUrl`
- `contractAddress`
- `deploymentStartBlock`
- `confirmations`

如果想按环境管理配置，参考：

- `backend/config/examples/anvil.example.json`
- `backend/config/examples/sepolia.example.json`

## 4. 启动后端

进入 `backend/`：

```bash
go mod tidy
```

PowerShell 示例：

```powershell
$env:DATABASE_URL = "root:password@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=true"
$env:CHAIN_CONFIG_PATH = "config/chain.testnet.example.json"
```

启动 API：

```bash
go run ./cmd/api
```

另一个终端启动 indexer：

```bash
go run ./cmd/indexer
```

健康检查：

```bash
curl http://127.0.0.1:8080/healthz
```

预期返回：

```json
{"status":"ok"}
```

## 5. 启动前端

在仓库根目录执行：

```bash
npx serve frontend
```

打开页面后：

1. 选择正确网络
2. 输入合约地址
3. 连接钱包
4. 测试 `create / fund / withdraw / refund`

## 6. 一轮自检

- `forge test -vv` 通过
- `GET /healthz` 返回正常
- indexer 开始推进 checkpoint
- MySQL 中出现 `campaigns` 等表
- 前端可以成功发起交易

## 常见问题

### RPC 配置错误

如果 indexer 或 API 无法连接 RPC，先检查：

- `rpcHttpUrl` 是否可达
- API Key 是否为真实值
- 宿主机能否直接访问该 RPC

### 数据库连接失败

检查：

- `DATABASE_URL` 是否正确
- MySQL 是否已启动
- `parseTime=true` 是否保留

### 交易成功但列表没有数据

这通常不是前端问题，而是：

- indexer 还没同步到对应区块
- `deploymentStartBlock` 设置错误
- `contractAddress` 配错

## 相关文档

- [docs/PROJECT_INTRO.md](PROJECT_INTRO.md)
- [docs/README.md](README.md)
- [docs/ARCHITECTURE.md](ARCHITECTURE.md)
- [docs/DEPLOYMENT.md](DEPLOYMENT.md)
- [docs/ONBOARDING.md](ONBOARDING.md)
- [backend/config/README.md](../backend/config/README.md)
