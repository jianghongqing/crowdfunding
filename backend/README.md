# Go Backend

这个目录是众筹项目的链下后端，职责是：

- 读取链上众筹合约状态
- 同步合约事件到 Postgres
- 提供查询 API 给前端

它**不**替用户签名交易，交易仍由钱包发起。

## 1) 配置

复制 `config/chain.testnet.example.json` 并填入：

- `rpcHttpUrl`
- `contractAddress`
- `deploymentStartBlock`

设置环境变量：

- `CHAIN_CONFIG_PATH`（可选，默认 `config/chain.testnet.example.json`）
- `DATABASE_URL`（必填）
- `API_ADDR`（可选，默认 `:8080`）

## 2) 生成合约 Go 绑定

本项目已放入：

- `contracts/crowdfund/CrowdFund.abi`
- `contracts/crowdfund/CrowdFund.bin`
- `contracts/crowdfund/crowdfund.go`

重新生成步骤：

1. 使用 Foundry 编译更新 `out/CrowdFund.sol/CrowdFund.json`
2. 执行：

```bash
./scripts/generate_bindings.sh
```

如果 `abigen` 不在默认路径，先设置：

```bash
export ABIGEN_BIN=/path/to/abigen
```

## 3) 数据库

执行 `migrations/001_init.sql` 初始化表结构。

## 4) 启动

安装依赖：

```bash
go mod tidy
```

启动索引器：

```bash
go run ./cmd/indexer
```

启动 API：

```bash
go run ./cmd/api
```

## 5) 接口

- `GET /healthz`
- `GET /campaigns`
- `GET /campaigns/{id}`
- `GET /campaigns/{id}/contributions/{address}`
