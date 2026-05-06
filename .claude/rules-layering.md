# 规则分层

为了让 Codex/Claude 更稳定地工作，建议把规则分成三层。

## 1. `AGENTS.md`

这里放强规则：

- 必须遵守的项目约束
- 不能做的事
- 测试与部署联动要求
- 架构边界

特点：

- 规则少但硬
- 尽量稳定
- 不放太多背景介绍

## 2. `CLAUDE.md`

这里放最短导航：

- 先读哪些文件
- 最核心的 3 到 5 条项目事实
- 最常见的联动规则

特点：

- 很短
- 适合快速建立上下文
- 不和 `AGENTS.md`、`.claude/context.md` 大量重复

## 3. `.claude/`

这里放辅助上下文：

- `context.md`：项目事实
- `development-rules.md`：协作方式
- `prompt-templates.md`：提需求模板
- `rules-layering.md`：这份分层说明

特点：

- 可扩展
- 面向团队协作
- 面向节省 token

## 推荐原则

- 强约束进 `AGENTS.md`
- 快速导航进 `CLAUDE.md`
- 细节和模板进 `.claude/`

这样通常比把所有内容都堆进一个文件更稳。
