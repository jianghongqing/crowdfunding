# Claude Quick Start

先不要扫全仓库，按这个顺序建立上下文：

1. `AGENTS.md`
2. `.claude/context.md`
3. `.claude/development-rules.md`
4. `docs/ARCHITECTURE.md`

## 只记住这 4 件事

- 合约是状态和资金真相源
- 前端直接发交易，后端不代签
- API 主要读库，必要时回链上兜底
- indexer 只做事件同步，不混入业务写操作

## 联动规则

- 改合约事件：同步检查 ABI、indexer、store、前端
- 改数据库结构：同步补 `backend/migrations/`
- 改部署方式：同步更新 `docker-compose.yml` 和 `docs/DEPLOYMENT.md`
- 改链配置路径：同步检查 `.env.example` 和 compose 挂载

## 安全底线

- 不提交真实密钥、生产配置、个人本地权限文件
- 不提交生成文件、二进制和异常备份文件
- 不在未确认影响范围前做大范围重构
