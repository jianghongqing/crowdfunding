# Go Backend

这个目录负责三件事：

- 监听链上 `CrowdFund` 合约事件并写入 MySQL
- 提供活动查询 API
- 在数据库没有命中时，回退到链上读取活动详情

## 目录说明

- `cmd/api`：HTTP API 入口
- `cmd/indexer`：链上事件索引器
- `internal/store`：MySQL 访问层
- `internal/chain`：合约读取和 RPC 连接
- `migrations`：数据库初始化与增量迁移

## 必要环境变量

- `DATABASE_URL`：MySQL DSN，必填
- `CHAIN_CONFIG_PATH`：链配置文件路径，默认 `config/chain.testnet.example.json`
- `API_ADDR`：API 监听地址，默认 `:8080`
- `DB_MAX_OPEN_CONNS`：连接池最大连接数，默认 `20`
- `DB_MAX_IDLE_CONNS`：连接池最大空闲连接数，默认 `10`
- `DB_CONN_MAX_LIFETIME`：连接最长生命周期，默认 `30m`
- `DB_CONN_MAX_IDLE_TIME`：连接最大空闲时间，默认 `5m`

示例：

```bash
export DATABASE_URL='root:password@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=true'
export CHAIN_CONFIG_PATH='config/chain.testnet.example.json'
export API_ADDR=':8080'
```

## 本地启动

1. 执行数据库初始化：

```bash
mysql -uroot -p crowdfunding < migrations/001_init.sql
mysql -uroot -p crowdfunding < migrations/002_add_campaign_status.sql
```

2. 启动索引器：

```bash
go run ./cmd/indexer
```

3. 启动 API：

```bash
go run ./cmd/api
```

## API

- `GET /healthz`
- `GET /campaigns?limit=20&offset=0`
- `GET /campaigns/{id}`
- `GET /campaigns/{id}/contributions/{address}`

`/campaigns` 支持分页，`limit` 范围为 `1-100`。

## 容器化部署

仓库根目录已经提供：

- `backend/Dockerfile`
- `docker-compose.yml`
- `.env.example`

启动方式：

```bash
cp .env.example .env
docker compose up -d --build
```

上线前请把 `backend/config/chain.testnet.example.json` 替换成真实链配置，尤其是：

- `rpcHttpUrl`
- `contractAddress`
- `deploymentStartBlock`
