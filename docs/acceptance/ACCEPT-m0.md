# M0 工程基线 —— 验收记录（Hermes）

日期：2026-09-14　基线：`562bd6b` → `1f55b91`（8 个原子 commit）

## 验收结论：✅ 通过

## 逐项核对（对照任务书 5 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | `bash scripts/check.sh` 四绿 | ✅ | Hermes 复跑：go vet/test（authstore/control/fsutil/upstream 全 ok）+ web lint(0 问题)/typecheck(0 错误)/test(39/39)/build(dist 162KB js + 9KB css) 全部通过 |
| 2 | 原子提交序列 test 先于 feat | ✅ | 6c058c8 工具链 → 9599a36/4c4af42 主题 → 29c8fea/9744e99 路由 → bd90be1/00b67ee 原语 → 932f126/1f55b91 页面拆分，4 组 RED→GREEN 配对完整 |
| 3 | 8 页面行为零回归 + dist 可构建 | ✅ | 前端 fetch URL 逐一与 `internal/control/server.go` 22 条路由核对一致（/api/session /login /logout /overview /accounts /accounts/:uid/actions/:a /credits /credits/:uid /models /automations(+PUT/run) /scheduler-tasks(+/detail) /mock-scheduler-tasks(CRUD) /activities /oauth/start /oauth/:id/poll）；`go build ./...` 通过（embed dist 正常） |
| 4 | 前端测试 ≥25 | ✅ | 39 个（6 个测试文件：app 2 / theme 3 / useHashRoute 5 / components 9 / pages 15 / routing 5） |
| 5 | token 浅色值与原 styles.css 逐字一致 + 暗色无新色相 | ✅ | 报告附 21 值 diff 为空；暗色 16 值全部为既有色相（薄荷绿/红/黄/灰）的暗色推导，主色 `#54c7ad` 原样保留 |

## Hermes 抽查（mock 盲区排查）

- services/api.ts 错误信封解析（`body.error` fallback）与后端 `{"error": ...}` 约定一致。
- 页面拆分后 App.tsx 仅剩壳+路由装配（55 行），9 页面全部独立文件，无行为漂移迹象。
- 覆盖率复跑：行 84.06%（达标），函数 42.25% —— 与执行体自述一致，无粉饰；已知留白（Automation/Scheduler 交互回调）记入 M1+ 契约测试范围。

## 执行体偏差说明核实（6 条全部合理）

1. ESLint `react-hooks/set-state-in-effect` 误报关闭——合理（标准 fetch-in-effect 模式，注释在案）。
2. 手写 30 行 hash 路由而非 react-router——合理（零依赖）。
3. 函数覆盖率 42%——如实披露，后续里程碑补。
4. Notice 组件提取——结构逐字一致。
5. 测试集中在 src/test/——无影响。
6. vitest 2.x（vite 5 peer 兼容）——合理，升级 vite 超范围。

## 遗留与后续

- `web/coverage/` 未入 .gitignore（工作区 untracked）→ M1 任务书补一条。
- 下一步：M1 契约终稿 + TS 类型 + MSW 双环境（任务书 `docs/milestones/m1-api-contract.md`）。
