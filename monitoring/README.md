# 监控与告警

这个目录用于沉淀运行期观测项和告警建议。

当前文件：

- `README.md`
- `checklist.md`
- `alerts-template.md`

## 最低监控集

- API 健康检查：`GET /healthz`
- Indexer 最新同步区块
- Indexer 落后区块数
- MySQL 连接数
- RPC 请求失败率

## 建议告警

- `healthz` 连续失败
- indexer 长时间不推进 checkpoint
- indexer 落后超过设定阈值
- MySQL 不可连接
- RPC 响应大量超时

## 推荐阈值

- indexer 60 秒未推进：告警
- indexer 落后超过 20 个确认后区块：告警
- API 连续 3 次健康检查失败：告警

## 后续可扩展

- Prometheus 指标
- Grafana 仪表盘
- 日志聚合
- 关键链上事件告警
