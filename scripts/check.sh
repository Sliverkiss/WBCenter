#!/usr/bin/env bash
# 四绿检查脚本：Go vet/test + 前端 lint/typecheck/test/build。
# 任一环节失败即退出非零。
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> go vet ./..."
go vet ./...

echo "==> go test ./..."
go test ./...

echo "==> web: lint"
(cd web && npm run lint)

echo "==> web: typecheck"
(cd web && npm run typecheck)

echo "==> web: test"
(cd web && npm test -- --run)

echo "==> web: build"
(cd web && npm run build)

echo "==> 全部通过"
