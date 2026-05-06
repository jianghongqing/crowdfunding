# 配置说明

当前后端运行依赖两类配置：

1. 环境变量
2. 链配置 JSON

## 环境变量

最关键的变量：

- `DATABASE_URL`
- `CHAIN_CONFIG_PATH`
- `API_ADDR`
- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`
- `DB_CONN_MAX_LIFETIME`
- `DB_CONN_MAX_IDLE_TIME`

## 链配置文件

当前代码默认示例文件：

- `backend/config/chain.testnet.example.json`

推荐把不同环境配置整理在：

- `backend/config/examples/`

当前已提供：

- `backend/config/examples/anvil.example.json`
- `backend/config/examples/sepolia.example.json`
- `backend/config/examples/production.template.json`

## 环境约定

建议团队统一使用下面三类命名：

- `local`
  - 本地开发或 Anvil
  - 追求快速联调
  - `confirmations` 可以较小
- `testnet`
  - Sepolia 等公开测试网
  - 用于演示、联调、预发布验证
- `production`
  - 面向正式环境
  - 不把真实配置提交进仓库
  - 仓库中只保留模板，不保留真实值

推荐做法：

- 本地运行：直接使用 `examples/anvil.example.json`
- 测试网联调：复制 `examples/sepolia.example.json`
- 生产部署：参考 `examples/production.template.json` 在运维环境单独生成真实文件

## 推荐实践

- 把样例文件提交到仓库
- 把真实环境配置保存在运维环境，不直接提交
- 通过 `.env` 的 `CHAIN_CONFIG_HOST_PATH` 指向实际挂载文件
- 修改 `contractAddress` 后，同步检查 `deploymentStartBlock`

## 字段说明

- `chainName`：链名称
- `chainId`：链 ID
- `rpcHttpUrl`：HTTP RPC 地址
- `rpcWsUrl`：WS RPC 地址，可选
- `contractAddress`：部署后的众筹合约地址
- `deploymentStartBlock`：indexer 从哪个区块开始扫
- `confirmations`：确认块数

## 相关文档

- [docs/START_FROM_ZERO.md](../../docs/START_FROM_ZERO.md)
- [docs/DEPLOYMENT.md](../../docs/DEPLOYMENT.md)
- [docs/ARCHITECTURE.md](../../docs/ARCHITECTURE.md)
