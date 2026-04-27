# Crowdfunding 项目从 0 运行指南

这份文档按当前仓库代码整理，后端数据库是 MySQL。

## Go 使用的框架

- HTTP 框架：`net/http` + `github.com/go-chi/chi/v5`
- 不是 Gin / Fiber
- 证据：`backend/internal/api/handlers.go` 中使用了 `chi.NewRouter()`

## 1. 环境准备

至少需要：

- Go 1.25+
- Node.js 18+
- MySQL 5.7+ / 8.0+
- Foundry（`forge` / `anvil`）
- 浏览器钱包（MetaMask）

## 2. 拉起合约（Foundry）

在项目根目录执行：

```bash
forge build
forge test -vv
```

如果你在 Windows 上通过 `bash` 使用 foundry，也可以：

```bash
bash -lc "~/.foundry/bin/forge build"
bash -lc "~/.foundry/bin/forge test -vv"
```

### 本地链部署（可选）

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

记录输出的合约地址和部署区块。

## 3. 初始化 MySQL

### 3.1 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS crowdfunding CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3.2 导入表结构

在 MySQL 客户端执行：

```sql
USE crowdfunding;
SOURCE backend/migrations/001_init.sql;
```

如果你的 MySQL 版本不支持 `CREATE INDEX IF NOT EXISTS`，可手动补索引：

```sql
CREATE INDEX idx_contributions_campaign_id ON contributions (campaign_id);
CREATE INDEX idx_refunds_campaign_id ON refunds (campaign_id);
CREATE INDEX idx_withdrawals_campaign_id ON withdrawals (campaign_id);
```

## 4. 配置链参数

编辑 `backend/config/chain.testnet.example.json`：

- `rpcHttpUrl`：可用 RPC（不能是占位 key）
- `contractAddress`：你部署出来的合约地址
- `deploymentStartBlock`：部署区块高度
- `confirmations`：确认块数（默认 5）

如果是 Sepolia，`chainId` 保持 `11155111` 即可。

## 5. 启动后端（API + Indexer）

先进入后端目录：

```bash
cd backend
go mod tidy
```

设置环境变量（PowerShell）：

```powershell
$env:DATABASE_URL = "root:password@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4"
$env:CHAIN_CONFIG_PATH = "config/chain.testnet.example.json"
```

启动 API：

```bash
go run ./cmd/api
```

启动 Indexer（另一个终端）：

```bash
go run ./cmd/indexer
```

### 健康检查

```bash
curl http://127.0.0.1:8080/healthz
```

期望输出：

```json
{"status":"ok"}
```

## 6. 启动前端

在项目根目录：

```bash
npx serve frontend
```

打开页面后：

1. 在「网络」里选 **Sepolia** 或 **本地 Anvil**（本地需先 `anvil`，链 ID `31337`）
2. 输入合约地址
3. 点击连接钱包（会尝试切换/添加对应网络）
4. 测试 `create / fund / withdraw / refund`

## 7. 一键检查清单

- `forge test -vv` 通过
- MySQL 里存在 5 张表：`campaigns` / `contributions` / `refunds` / `withdrawals` / `indexer_checkpoints`
- `GET /healthz` 返回 `ok`
- indexer 不再报 `401 Unauthorized`
- 前端可正常发交易

## 常见问题

### 1) `401 Unauthorized: invalid project id`

`rpcHttpUrl` 仍是占位的 Infura key，换成真实可用 RPC。

### 2) `go run` 报找不到可执行文件（Windows）

检查 Go 目标平台是否被改过：

```bash
go env GOOS GOARCH
```

如果不是 `windows/amd64`，可在当前终端修正：

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
```

### 3) MySQL 连接串里密码含 `@`

优先使用标准 MySQL DSN 形式：

`user:pass@tcp(host:port)/dbname?charset=utf8mb4`

必要时对特殊字符做转义或改用不含特殊字符的密码。

### 4) `wsl --update` / `wsl --install` 报 `0x80072f7d`（安全频道）

浏览器能开 GitHub，但 WSL 在线更新失败时：

1. **以管理员**运行仓库内脚本：`scripts/fix-wsl.ps1`（从 GitHub 安装官方 `wsl.*.x64.msi`）。
2. 若短期不需要 WSL：运行 `scripts/install-foundry-windows.ps1` 使用 **Windows 原生** `forge` / `anvil`。
