# 部署说明

本文档描述当前项目推荐的部署方式、配置项和上线检查项。

## 部署拆分

推荐拆分为 4 个服务：

1. `frontend`
2. `api`
3. `indexer`
4. `mysql`

原因：

- 前端是静态资源，适合单独托管
- API 是读服务，负责查询响应
- indexer 是后台同步任务，负责持续追块
- MySQL 是持久化层，应该独立管理

## 推荐方式

- `frontend`：Nginx、对象存储 + CDN、Vercel 等静态托管
- `api`：容器单独部署，放在反向代理后
- `indexer`：单独容器部署，不与 API 混跑
- `mysql`：托管数据库或专用主机

## 配置文件

链配置由宿主机文件挂载进入容器，推荐通过 `.env` 里的 `CHAIN_CONFIG_HOST_PATH` 指定。

示例：

```env
CHAIN_CONFIG_HOST_PATH=./backend/config/chain.testnet.example.json
```

容器内统一读取：

```text
/app/config/chain.json
```

## 必要环境变量

- `MYSQL_ROOT_PASSWORD`
- `MYSQL_DATABASE`
- `MYSQL_PORT`
- `API_PORT`
- `DATABASE_URL`
- `CHAIN_CONFIG_HOST_PATH`
- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`
- `DB_CONN_MAX_LIFETIME`
- `DB_CONN_MAX_IDLE_TIME`

## 启动方式

```bash
cp .env.example .env
docker compose up -d --build
```

## 发布顺序

1. 准备链配置文件
2. 启动 MySQL
3. 确认 migrations 已执行
4. 启动 indexer
5. 启动 API
6. 检查 `GET /healthz`
7. 发布前端

## 上线前检查

- `CHAIN_CONFIG_HOST_PATH` 指向真实配置，而不是占位样例
- `contractAddress` 与当前部署合约一致
- `deploymentStartBlock` 与真实部署区块一致
- API 容器可以连通 RPC
- Indexer 能够推进 checkpoint
- MySQL 只对内网开放
- 敏感变量不写死在前端代码里

## 运行期建议

- 只暴露前端和 API 对外入口
- 将 RPC Key、数据库凭据放入环境变量或密钥管理服务
- 为 API 和 indexer 配置自动重启策略
- 每日备份 MySQL
- 监控以下指标：
  - API `healthz`
  - indexer 最新同步区块
  - MySQL 连接数
  - RPC 错误率

## 相关文档

- [docs/PROJECT_INTRO.md](PROJECT_INTRO.md)
- [docs/README.md](README.md)
- [docs/ARCHITECTURE.md](ARCHITECTURE.md)
- [docs/START_FROM_ZERO.md](START_FROM_ZERO.md)
- [docs/ONBOARDING.md](ONBOARDING.md)
- [backend/config/README.md](../backend/config/README.md)
- [monitoring/README.md](../monitoring/README.md)
