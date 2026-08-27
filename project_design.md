# 项目设计

## 业务场景

- 业务主体：发布协调员、审阅人、批准人和值班负责人。
- 用户动作：创建发布包、补充说明、提交审阅、批准发布、执行发布、回滚已发布版本。
- 状态流转：draft -> submitted -> reviewed -> approved -> published -> rolled_back。
- 持久化边界：所有状态变化先写入 JSONL journal，再生成快照；重启后从 journal 回放并恢复快照。
- 可观察结果：命令行输出当前状态、审阅记录、批准信息、发布摘要和回滚原因。

## 功能映射

| 文档功能 | Go 入口 | 下游影响 | 验证路径 |
|---|---|---|---|
| 提交发布包 | `internal/service` | 写 journal、更新快照 | 单元测试、CLI smoke |
| 审阅与批准 | `internal/service` | 状态推进、审计记录 | 状态机测试 |
| 发布与回滚 | `internal/service` + `internal/storage` | 生成报告、恢复检查点 | 恢复测试 |
| 导出报表 | `internal/storage` | JSON 报表输出 | 导出测试 |

