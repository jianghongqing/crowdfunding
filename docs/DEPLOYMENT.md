# 部署建议

这个项目更适合拆成 4 个角色部署：

1. `frontend`
2. `api`
3. `indexer`
4. `mysql`

## 推荐方案

- 前端：静态托管到 Nginx、Vercel 或 OSS + CDN
- API：单独容器，挂在反向代理后面
- Indexer：独立容器，避免和 API 互相影响
- MySQL：托管数据库或独立云主机，不建议长期和应用混部

## 为什么这样拆

- API 和索引器负载模型不同：一个偏读请求，一个偏链上同步
- 索引器重启、追块、补块时，不会拖慢 API
- 前端是静态资源，应该走 CDN，减少源站压力

## 当前仓库已落地的部署优化

- API 已加请求超时、中间件、优雅关闭
- 索引器改成按区块批次同步，减少大区间 `FilterLogs` 失败概率
- MySQL 连接池支持环境变量配置
- `campaigns` 接口支持分页，避免列表无限长
- 活动状态直接落库，前端列表读取更轻

## 线上建议

- Nginx 反代 `api`，只暴露 `80/443`
- `mysql` 只允许内网访问
- `DATABASE_URL`、RPC Key、合约地址都放进环境变量或密钥管理
- API 和 Indexer 分别配置重启策略
- 给 MySQL 做每日备份
- 监控最少覆盖：容器存活、API `healthz`、索引器最新区块、MySQL 连接数

## 发布顺序

1. 先部署 MySQL
2. 执行 `migrations/001_init.sql` 和 `migrations/002_add_campaign_status.sql`
3. 部署 Indexer 并确认开始追块
4. 部署 API 并检查 `GET /healthz`
5. 最后部署前端

## 后续还能继续优化的点

- 给 API 增加 CORS 白名单和速率限制
- 给 Indexer 增加 Prometheus 指标
- 引入消息队列，分离重算与回填任务
- 前端把 API 地址抽成配置，而不是写死在页面脚本里
