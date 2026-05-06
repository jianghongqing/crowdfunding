# 发布前安全检查清单

发布前建议逐项确认：

## 合约

- [ ] 当前部署的合约地址已经记录
- [ ] 当前部署区块已经记录
- [ ] 合约 ABI 与前后端使用版本一致
- [ ] 关键流程已回归测试：
  - [ ] `createCampaign`
  - [ ] `fund`
  - [ ] `withdraw`
  - [ ] `refund`
- [ ] 已确认没有引入新的管理员权限、暂停逻辑或升级逻辑

## 配置

- [ ] `contractAddress` 指向正确部署
- [ ] `deploymentStartBlock` 与真实部署块一致
- [ ] `confirmations` 配置符合目标网络特性
- [ ] 没有把真实私钥、RPC Key、数据库密码提交到仓库
- [ ] `CHAIN_CONFIG_HOST_PATH` 指向真实环境配置，不是示例占位文件

## 后端

- [ ] API 能正常启动
- [ ] Indexer 能正常启动
- [ ] `GET /healthz` 返回正常
- [ ] indexer 能推进 checkpoint
- [ ] MySQL 迁移已执行完毕

## 前端

- [ ] 前端显示的链和合约地址正确
- [ ] 钱包连接流程正常
- [ ] API 地址没有硬编码成错误环境
- [ ] 活动列表和详情能正常展示

## 仓库卫生

- [ ] 没有提交 `.exe`、`.test`、`.go.数字` 这类本地产物
- [ ] 没有提交 `.claude/settings.local.json`
- [ ] 部署方式变更已同步更新文档
- [ ] 数据库结构变更已同步补 migration
