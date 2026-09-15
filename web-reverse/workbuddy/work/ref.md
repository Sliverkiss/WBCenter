# Growth task implementation references

- Upstream runner: https://github.com/Sliverkiss/workbuddy2api/blob/master/scripts/task_runner.py
- Community panel: https://github.com/linguo2625469/workbuddy2api-panel
- Alternative reward/travel implementation: https://github.com/88lin/workbuddy-auto-signin
- Feature context: https://github.com/Sliverkiss/workbuddy2api/issues/61

The control center mounts the upstream `scripts/` directory read-only and
executes `task_runner.py` without a shell. Account IDs come from the local auth
store and task codes are restricted to a compiled allowlist.
