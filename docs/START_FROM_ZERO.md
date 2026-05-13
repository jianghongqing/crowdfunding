# 从零启动（本地开发完整指南）

本文档面向第一次在本地跑起整个项目的开发者，逐步说明每条命令在哪个目录执行、每个配置从哪里获取、填到哪里。

## 依赖准备

| 工具 | 最低版本 | 用途 |
|------|---------|------|
| Go | 1.25+ | 编译后端服务 |
| Foundry (forge, anvil) | latest | 合约编译/测试/部署/本地链 |
| Node.js | 18+ | 前端静态服务 |
| MySQL | 8+ | 后端数据存储 |
| 浏览器钱包 (MetaMask) | — | 与合约交互 |

---

## 1. 编译与测试合约

**执行目录：** 仓库根目录 `g:\web3\crowdfunding`

```bash
forge build
forge test -vv
```

确认全部通过后继续。

---

## 2. 启动本地链并部署合约

### 2.1 启动 Anvil（终端 A）

**执行目录：** 任意（Anvil 是全局命令）

```bash
anvil
```

启动后终端会输出 10 个测试账户，格式如下：

```
Available Accounts
==================
(0) 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10000 ETH)
...

Private Keys
==================
(0) 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
...
```

> **记下第 (0) 个 Private Key**，下一步部署要用。

### 2.2 部署合约（终端 B）

**执行目录：** 仓库根目录 `g:\web3\crowdfunding`

```bash
forge script script/DeployCrowdFund.s.sol:DeployCrowdFundScript ^
  --rpc-url http://127.0.0.1:8545 ^
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 ^
  --broadcast
```

> 上面的 private key 是 Anvil 默认第 0 个账户的私钥，仅限本地测试使用。

部署成功后终端会输出类似：

```
== Return ==
crowdFund: address 0x5FbDB2315678afecb367f032d93F642f64180aa3
...
Block Number: 1
```

**记录两个值：**

| 信息 | 示例值 | 后续填入位置 |
|------|--------|-------------|
| 合约地址 | `0x5FbDB2315678afecb367f032d93F642f64180aa3` | 链配置 JSON 的 `contractAddress` |
| 部署区块号 | `1` | 链配置 JSON 的 `deploymentStartBlock` |

---

## 3. 初始化 MySQL

### 3.1 创建数据库

用 MySQL 客户端登录后执行：

```sql
CREATE DATABASE IF NOT EXISTS crowdfunding CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3.2 执行迁移脚本

```sql
USE crowdfunding;
SOURCE backend/migrations/001_init.sql;
SOURCE backend/migrations/002_add_campaign_status.sql;
```

> `SOURCE` 的路径是相对于你启动 MySQL 客户端时所在的目录。如果在仓库根目录启动 `mysql -u root -p`，上面的路径即可直接使用。

或者直接在命令行执行（**执行目录：** 仓库根目录）：

```bash
mysql -u root -p crowdfunding < backend/migrations/001_init.sql
mysql -u root -p crowdfunding < backend/migrations/002_add_campaign_status.sql
```

---

## 4. 准备链配置文件

### 4.1 文件位置与用途

后端服务启动时通过环境变量 `CHAIN_CONFIG_PATH` 读取一个 JSON 文件，默认值为：

```
config/chain.testnet.example.json
```

（相对于 `backend/` 目录）

该文件的完整路径是：

```
g:\web3\crowdfunding\backend\config\chain.testnet.example.json
```

### 4.2 文件内容与字段说明

打开上面的文件，内容结构如下：

```json
{
  "chainName": "anvil",
  "chainId": 31337,
  "rpcHttpUrl": "http://127.0.0.1:8545",
  "rpcWsUrl": "ws://127.0.0.1:8545",
  "contractAddress": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "deploymentStartBlock": 1,
  "confirmations": 1
}
```

**逐字段说明：**

| 字段 | 含义 | 值从哪里来 | 是否必填 |
|------|------|-----------|---------|
| `chainName` | 链名称标识，用于日志展示 | 本地链填 `anvil`；测试网填 `sepolia` | 建议填写 |
| `chainId` | 链 ID | Anvil 默认 `31337`；Sepolia 为 `11155111` | 建议填写 |
| `rpcHttpUrl` | HTTP RPC 端点 | Anvil 本地为 `http://127.0.0.1:8545`；Sepolia 从 Infura/Alchemy 获取 | **必填** |
| `rpcWsUrl` | WebSocket RPC 端点 | 同上对应的 WS 地址；当前后端未使用，可留空 | 可选 |
| `contractAddress` | 已部署的合约地址 | **第 2 步部署时输出的合约地址** | **必填** |
| `deploymentStartBlock` | indexer 从哪个区块开始扫描事件 | **第 2 步部署时输出的 Block Number** | **必填** |
| `confirmations` | 确认块数，indexer 只扫 `最新块 - confirmations` 之前的数据 | 本地链设 `1`；测试网建议 `5`；主网建议 `12` | 可选（默认 5） |

### 4.3 根据你的部署结果修改

如果你按照第 2 步部署后得到的地址是 `0xABC...123`、区块号是 `5`，则修改为：

```json
{
  "chainName": "anvil",
  "chainId": 31337,
  "rpcHttpUrl": "http://127.0.0.1:8545",
  "rpcWsUrl": "ws://127.0.0.1:8545",
  "contractAddress": "0xABC...123",
  "deploymentStartBlock": 5,
  "confirmations": 1
}
```

### 4.4 其他参考模板

| 文件路径 | 适用场景 |
|---------|---------|
| `backend/config/examples/anvil.example.json` | 本地 Anvil 链 |
| `backend/config/examples/sepolia.example.json` | Sepolia 测试网 |
| `backend/config/examples/production.template.json` | 生产环境模板 |

> 如果使用 Sepolia 测试网，需要去 [Infura](https://infura.io) 或 [Alchemy](https://alchemy.com) 注册免费账号，创建项目获取 RPC URL，替换 `<your-project-id>` 占位符。

---

## 5. 启动后端

### 5.1 安装 Go 依赖

**执行目录：** `g:\web3\crowdfunding\backend`

```bash
cd backend
go mod tidy
```

### 5.2 设置环境变量

在当前终端设置环境变量（PowerShell）：

```powershell
$env:DATABASE_URL = "root:your_password@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=true"
$env:CHAIN_CONFIG_PATH = "config/chain.testnet.example.json"
```

> **注意：**
> - `your_password` 替换成你的 MySQL root 密码
> - `CHAIN_CONFIG_PATH` 是**相对于 `backend/` 目录**的路径
> - 如果你把配置复制到了别处，使用绝对路径也可以

Linux/Mac 等效命令：

```bash
export DATABASE_URL="root:your_password@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=true"
export CHAIN_CONFIG_PATH="config/chain.testnet.example.json"
```

### 5.3 启动 API 服务（终端 C）

**执行目录：** `g:\web3\crowdfunding\backend`

```bash
go run ./cmd/api
```

成功输出：

```
api listening on :8080
```

### 5.4 启动 Indexer 服务（终端 D）

**执行目录：** `g:\web3\crowdfunding\backend`

（同样需要设置 5.2 中的环境变量）

```bash
go run ./cmd/indexer
```

成功输出：

```
indexer started for contract 0x5FbDB2315678afecb367f032d93F642f64180aa3
```

### 5.5 健康检查

```bash
curl http://127.0.0.1:8080/healthz
```

预期返回：

```json
{"status":"ok"}
```

> PowerShell 中如果没有 curl，可用：`Invoke-RestMethod http://127.0.0.1:8080/healthz`

---

## 6. 启动前端

**执行目录：** 仓库根目录 `g:\web3\crowdfunding`

```bash
npx serve frontend
```

或者用任意静态文件服务器（如 Python）：

```bash
cd frontend
python -m http.server 3000
```

> 前端默认使用同源 API 路径，最适合通过 Docker Compose 的 nginx 前端容器访问。若使用 `npx serve frontend` 单独启动静态服务，可在浏览器控制台执行 `localStorage.setItem('CROWDFUND_API_BASE', 'http://127.0.0.1:8080')` 后刷新页面。

### 6.1 页面操作流程

1. 打开浏览器访问 `http://localhost:3000`（或 `npx serve` 输出的地址）
2. 在 Header 右侧选择网络（Sepolia 测试网 / 本地 Anvil）
3. 点击「连接钱包」，MetaMask 弹窗确认（合约地址已从 API `/config` 自动加载）
4. 使用 Tab 切换查看活动列表 / 发起活动 / 快捷操作
5. 测试功能：`create / fund / withdraw / refund`

### 6.2 MetaMask 添加本地 Anvil 网络

如果使用本地链，需在 MetaMask 手动添加网络：

| 设置项 | 值 |
|--------|-----|
| 网络名称 | Anvil Local |
| RPC URL | http://127.0.0.1:8545 |
| Chain ID | 31337 |
| 货币符号 | ETH |

然后导入 Anvil 测试账户私钥以获得测试 ETH。

---

## 7. 一轮自检清单

| 检查项 | 预期结果 |
|--------|---------|
| `forge test -vv` | 全部 PASS |
| `GET /healthz` | 返回 `{"status":"ok"}` |
| indexer 终端日志 | 持续输出 checkpoint 推进信息 |
| MySQL `SELECT * FROM campaigns` | 表已存在（可能暂无数据） |
| 前端页面创建活动后刷新列表 | 几秒后数据出现 |

---

### 6.3 使用 Docker Compose 一键启动前端

如果不想手动启动静态服务，可以使用 Docker Compose（前提是后端和 MySQL 也在容器中运行）：

```bash
docker compose up frontend -d
```

此模式下前端由 nginx 伺服，且自动反向代理 API 请求，`app.js` 默认同源 API 路径即可工作。

---

## 完整终端总览

本地开发至少需要 4 个终端窗口：

| 终端 | 目录 | 命令 | 用途 |
|------|------|------|------|
| A | 任意 | `anvil` | 本地 EVM 链 |
| B | 仓库根目录 | `forge script ...` | 部署合约（一次性） |
| C | `backend/` | `go run ./cmd/api` | API 服务 |
| D | `backend/` | `go run ./cmd/indexer` | 事件索引服务 |
| E | 仓库根目录 | `npx serve frontend` | 前端静态服务 |

> 替代方案：使用 `docker compose up -d --build` 一次性启动全部服务（包含 nginx 前端代理）。

---

## 常见问题

### RPC 配置错误

如果 indexer 或 API 报 `dial rpc` 错误：

- 确认 Anvil 是否在运行
- 确认 `rpcHttpUrl` 是否为 `http://127.0.0.1:8545`
- 如果用 Sepolia：确认 Infura/Alchemy API Key 是否有效

### 数据库连接失败

- 确认 `DATABASE_URL` 格式正确：`user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true`
- 确认 MySQL 已启动：`mysql -u root -p -e "SELECT 1"`
- `parseTime=true` 必须保留，否则时间字段解析会失败

### 交易成功但前端列表没有数据

这通常是 indexer 问题：

- indexer 还没同步到对应区块（等几秒）
- `deploymentStartBlock` 设置比合约部署区块大（检查配置）
- `contractAddress` 不一致（检查配置 vs 实际部署地址）

### PowerShell 环境变量只在当前终端生效

每个 PowerShell 窗口需要**单独设置**环境变量。终端 C 和 D 都需要执行 5.2 中的 `$env:` 命令。

---

## 相关文档

- [docs/PROJECT_INTRO.md](PROJECT_INTRO.md)
- [docs/README.md](README.md)
- [docs/ARCHITECTURE.md](ARCHITECTURE.md)
- [docs/DEPLOYMENT.md](DEPLOYMENT.md)
- [docs/PRODUCTION_DEPLOY.md](PRODUCTION_DEPLOY.md)
- [docs/ONBOARDING.md](ONBOARDING.md)
- [backend/config/README.md](../backend/config/README.md)
