# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

compact 在关闭「支撑事件已过期」的 incident 时，如果这次状态切换被拒绝，命令仍然返回成功，并且 closed_incidents 计数照常增加，输出的 incident 状态却还是 open。请修复这条 compaction 路径：状态切换被拒绝时必须把失败作为错误返回，不要静默计入关闭数量，同时不影响正常关闭与保留逻辑，并保证全量测试通过。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-46
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-46.git
- parent SHA：e2340b13e539801736a14b510d10ac6c56e3e959

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-46.git bug-repro
cd bug-repro
git checkout --detach e2340b13e539801736a14b510d10ac6c56e3e959
go test ./internal/retention -run "^TestCompactPropagatesFailedIncidentClosure$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/retention -run "^TestCompactPropagatesFailedIncidentClosure$" -count=1 -v
=== RUN   TestCompactPropagatesFailedIncidentClosure
    closure_regression_test.go:23: Compact reported success after a rejected closure: {BeforeEvents:0 AfterEvents:0 ExpiredEvents:0 DuplicateEvents:0 BeforeIncidents:1 AfterIncidents:1 ExpiredIncidents:0 ClosedIncidents:1 Cutoff:2025-01-09 00:00:00 +0000 UTC}
--- FAIL: TestCompactPropagatesFailedIncidentClosure (0.00s)
FAIL
FAIL	LogPilot/internal/retention	0.002s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/retention -run "^TestCompactPropagatesFailedIncidentClosure$" -count=1 -v
=== RUN   TestCompactPropagatesFailedIncidentClosure
    closure_regression_test.go:23: Compact reported success after a rejected closure: {BeforeEvents:0 AfterEvents:0 ExpiredEvents:0 DuplicateEvents:0 BeforeIncidents:1 AfterIncidents:1 ExpiredIncidents:0 ClosedIncidents:1 Cutoff:2025-01-09 00:00:00 +0000 UTC}
--- FAIL: TestCompactPropagatesFailedIncidentClosure (0.01s)
FAIL
FAIL	LogPilot/internal/retention	0.121s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

关闭 incident 被拒绝时 Compact 返回错误且不虚增 closed_incidents；正常关闭、过期与去重路径不回归；双架构定向、全量、build/vet 通过。
