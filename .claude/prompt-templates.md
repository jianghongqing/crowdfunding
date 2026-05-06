# Codex 提示模板

下面这些模板用于减少重复解释上下文、降低 token 消耗，并让 Codex 更稳定地先读对文件再动手。

## 1. 改代码前的通用模板

```text
先遵守 AGENTS.md，并先阅读：
1. .claude/context.md
2. .claude/development-rules.md
3. docs/ARCHITECTURE.md

任务：
<在这里写你的具体需求>

要求：
- 先理解现状再修改
- 只做必要的最小改动
- 不要动无关文件
- 如果涉及 migration / deployment / ABI / indexer 联动，请同步处理
```

## 2. 后端改动模板

```text
先遵守 AGENTS.md，并先阅读：
1. .claude/context.md
2. .claude/development-rules.md
3. backend/config/README.md
4. docs/ARCHITECTURE.md

我要修改 Go 后端：
<具体需求>

请重点检查：
- API 和 indexer 兼容性
- 是否需要 migration
- 是否需要更新 docs/DEPLOYMENT.md
- 是否需要补测试
```

## 3. 合约改动模板

```text
先遵守 AGENTS.md，并先阅读：
1. .claude/context.md
2. docs/PROJECT_INTRO.md
3. docs/ARCHITECTURE.md

我要修改智能合约：
<具体需求>

请重点检查：
- 是否影响 ABI
- 是否影响前端调用
- 是否影响 indexer/store
- 是否需要补 Foundry 测试
```

## 4. 前端改动模板

```text
先遵守 AGENTS.md，并先阅读：
1. .claude/context.md
2. .claude/development-rules.md
3. docs/ARCHITECTURE.md

我要修改 frontend：
<具体需求>

要求：
- 保持前端直连钱包，不把交易职责挪到后端
- 不随意改动现有页面结构
- 如果改接口，检查前后端兼容性
```

## 5. 文档整理模板

```text
先遵守 AGENTS.md，并先阅读：
1. .claude/context.md
2. docs/README.md

我要整理文档：
<具体需求>

要求：
- 优先补入口和交叉链接
- 不重复堆砌相同内容
- 让新人能按阅读顺序快速建立上下文
```

## 6. 代码评审模板

```text
请按 review 模式处理，先看：
1. AGENTS.md
2. security/review-checklist.md
3. .claude/context.md

评审范围：
<提交、文件或功能范围>

要求：
- 先给 findings，按严重程度排序
- 重点看 bug、回归风险、权限风险、配置风险
- 总结放后面
```

## 7. 发布前检查模板

```text
先遵守 AGENTS.md，并先阅读：
1. security/release-checklist.md
2. monitoring/checklist.md
3. docs/DEPLOYMENT.md

我准备发布这一版：
<发布范围>

请帮我检查：
- 是否有遗漏的配置项
- 是否有文档未同步
- 是否有明显的安全或发布风险
```

## 8. 最省 token 的写法建议

- 直接点名要读的 2 到 4 个文件
- 明确“只改相关文件，不扫全仓库”
- 明确“先理解现状，再最小改动”
- 明确是否要检查 migration、ABI、部署文档、测试

## 9. 不推荐的写法

- “帮我看看整个项目”
- “把这个项目重构一下”
- “你自己先看看哪里有问题”

这类提示范围太大，容易：

- 浪费 token
- 扫描过多无关文件
- 做出超出预期的大改动
