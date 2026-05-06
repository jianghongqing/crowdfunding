# 新成员上手指南

这份文档给第一次接手本仓库的开发者使用，目标是用最少路径建立正确上下文。

## 1. 先理解项目是什么

先读：

1. [文档索引](README.md)
2. [项目介绍](PROJECT_INTRO.md)
3. [架构说明](ARCHITECTURE.md)

你需要先建立一个基本认识：

- 合约是资金与状态真相源
- 前端直接发交易
- API 主要负责查询
- indexer 负责把链上事件同步到数据库

## 2. 再理解仓库重点

接着读：

1. [CLAUDE 快速入口](../CLAUDE.md)
2. [.claude/context.md](../.claude/context.md)
3. [.claude/development-rules.md](../.claude/development-rules.md)

这些文件的作用是减少重复扫仓库和重复解释上下文的 token 消耗。

## 3. 看核心代码入口

推荐顺序：

1. `src/CrowdFund.sol`
2. `test/CrowdFund.t.sol`
3. `script/DeployCrowdFund.s.sol`
4. `backend/cmd/api/main.go`
5. `backend/cmd/indexer/main.go`
6. `backend/internal/indexer/sync.go`
7. `backend/internal/store/`
8. `backend/internal/api/`

## 4. 跑起项目

先读：

1. [从零启动](START_FROM_ZERO.md)
2. [配置说明](../backend/config/README.md)

你要特别确认：

- `DATABASE_URL`
- `CHAIN_CONFIG_PATH`
- `contractAddress`
- `deploymentStartBlock`

## 5. 提交改动前要知道什么

先读：

1. [权限矩阵](../security/permissions-matrix.md)
2. [已知风险](../security/known-risks.md)
3. [代码评审安全清单](../security/review-checklist.md)

如果涉及部署或发布，再看：

1. [部署说明](DEPLOYMENT.md)
2. [发布前安全检查清单](../security/release-checklist.md)

## 6. 运行维护相关

先读：

1. [监控说明](../monitoring/README.md)
2. [运行巡检清单](../monitoring/checklist.md)
3. [告警模板](../monitoring/alerts-template.md)

## 7. 最容易犯错的地方

- 把前端交易职责改进后端
- 改合约事件后忘了同步 indexer
- 改数据库结构后忘了补 migration
- 把示例配置误当成生产配置
- 把 `.claude/settings.local.json` 提交进仓库

## 8. 一个建议的首日路径

第一天最适合做的事情：

1. 把文档看完
2. 本地跑通 `forge test`
3. 本地跑通 API 和 indexer
4. 用前端完成一次创建活动和捐款
5. 再开始修改业务代码
