# WorkBuddy growth tasks

Last verified: 2026-09-15

## Sources

- `Sliverkiss/workbuddy2api/scripts/task_runner.py`
- `linguo2625469/workbuddy2api-panel/internal/panel/autotask.go`
- `linguo2625469/workbuddy2api-panel/internal/upstream/tasks.go`

These are reverse-engineered, non-public WorkBuddy endpoints. Treat the schema
as unstable and always re-read task state after a write.

## Confirmed workflow

1. List: `GET /v2/activity/growth/tasks`.
2. Accept: `POST /activity/growth/tasks/accept` with
   `{"task_codes":["..."]}`.
3. Emit only a known task event through the appropriate `/v2/report` host.
4. Re-read the task list and verify `progress.current >= progress.target`.
5. Claim: `POST /activity/growth/tasks/{task_code}/claim` with no body.
6. A claim response containing `already_claimed` is idempotent success.

The older `/v2/activity/growth/tasks/reward/claim` route is not used by the
integrated runner.

## Safety boundary

Only task codes in the runner's reviewed mapping may be automated. Unknown
tasks and `Expert_Philanthropy` are shown as manual because the latter implies
a real donation side effect. `black_cat` is time-window constrained by the
runner. HTTP 200 from `/v2/report` is not completion evidence; the re-read is.
