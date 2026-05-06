# 告警模板

可以基于下面几类事件建立告警：

## API 健康检查失败

- 名称：`api_healthz_failed`
- 触发条件：连续 3 次请求 `/healthz` 失败
- 建议动作：检查 API 进程、日志、数据库连接、RPC 可达性

## Indexer 停滞

- 名称：`indexer_checkpoint_stalled`
- 触发条件：60 秒内 checkpoint 无推进
- 建议动作：检查 RPC、数据库、最近部署变更

## Indexer 落后过多

- 名称：`indexer_block_lag_high`
- 触发条件：落后安全头超过 20 个区块
- 建议动作：检查节点性能、FilterLogs 错误、数据库写入情况

## MySQL 不可用

- 名称：`mysql_unreachable`
- 触发条件：应用无法建立数据库连接
- 建议动作：检查数据库实例、网络、凭据与连接池设置
