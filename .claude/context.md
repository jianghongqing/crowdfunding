# 项目上下文

## 项目类型

这是一个众筹 DApp，采用：

- 合约保存业务真相和资金
- 前端直连钱包发交易
- Go API 提供查询接口
- Go indexer 把链上事件同步到 MySQL

## 关键入口

- 合约：`src/CrowdFund.sol`
- 合约测试：`test/CrowdFund.t.sol`
- 部署脚本：`script/DeployCrowdFund.s.sol`
- API 入口：`backend/cmd/api/main.go`
- Indexer 入口：`backend/cmd/indexer/main.go`
- 配置：`backend/config/`
- 部署文档：`docs/DEPLOYMENT.md`

## 先读什么

建议阅读顺序：

1. `docs/PROJECT_INTRO.md`
2. `docs/ARCHITECTURE.md`
3. `src/CrowdFund.sol`
4. `test/CrowdFund.t.sol`
5. `backend/internal/indexer/sync.go`
6. `backend/internal/store/`
7. `backend/internal/api/`

## 代码边界

- `frontend/`
  - 钱包交互和页面展示
- `src/`
  - 业务规则与资金流转
- `backend/internal/api`
  - 查询接口
- `backend/internal/indexer`
  - 事件同步
- `backend/internal/store`
  - MySQL 读写
- `backend/internal/chain`
  - 链上只读访问和状态推导

## 常见联动

- 改合约事件：检查后端绑定、indexer、store、前端 ABI
- 改数据库字段：补 migration
- 改后端行为：同时考虑 API 和 indexer 兼容性
- 改部署或配置：同步更新文档与 compose

## 节省 token 的做法

- 先读入口和文档，不先扫全仓库
- 用 `rg` 定位，再按需展开文件
- 先看 `docs/`、`.claude/`、`backend/config/README.md`
- 前端文件只在明确需要时再深入展开
