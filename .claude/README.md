# `.claude/` 怎么用

这个目录不是强规则来源，强规则以 `AGENTS.md` 为准。

它的用途是：

- 给 Claude/Codex 提供低成本项目上下文
- 帮新成员快速找到入口
- 减少每次重复解释仓库背景的 token 消耗

## 文件分工

- `context.md`
  - 高密度项目事实、入口、联动关系
- `development-rules.md`
  - 协作方式和改动前自查
- `prompt-templates.md`
  - 给 Codex 提需求时可直接复用的模板
- `rules-layering.md`
  - `AGENTS.md` / `CLAUDE.md` / `.claude/` 的分层说明
- `settings.local.example.json`
  - 本地权限配置示例
- `settings.local.json`
  - 个人本地配置，不提交

## 推荐用法

在提需求时直接点名：

```text
先遵守 AGENTS.md，并先阅读：
1. .claude/context.md
2. .claude/development-rules.md
3. docs/ARCHITECTURE.md
```

这样通常比等待模型自己发现 `.claude/` 更稳定。
