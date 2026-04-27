# AGENTS.md

## Project Overview

This repository is a crowdfunding DApp with three main parts:

- `src/`, `script/`, `test/`: Solidity contract, Foundry scripts, and Foundry tests
- `backend/`: Go services, including an API server and an on-chain event indexer
- `frontend/`: static frontend for wallet interaction and demo flows

Core contract:

- `src/CrowdFund.sol`

Core backend entrypoints:

- `backend/cmd/api/main.go`
- `backend/cmd/indexer/main.go`

Important backend modules:

- `backend/internal/api`
- `backend/internal/indexer`
- `backend/internal/store`
- `backend/internal/chain`
- `backend/internal/config`

Deployment assets:

- `backend/Dockerfile`
- `docker-compose.yml`
- `docs/DEPLOYMENT.md`

## Architecture Notes

- The frontend sends wallet-signed transactions directly to the chain.
- The Go `indexer` syncs contract events into MySQL.
- The Go `api` reads mostly from MySQL and may fall back to chain reads for some detail queries.
- The backend assumes a chain config file plus `DATABASE_URL`.

## Working Rules For Codex

- Prefer small, local changes over broad refactors.
- Preserve the current architecture: contract, indexer, API, and static frontend remain separate concerns.
- Do not rewrite frontend structure unless the task explicitly asks for frontend work.
- Do not regenerate contract bindings unless contract ABI or bytecode actually changed.
- When changing backend behavior, keep API and indexer compatibility in mind together.
- When changing database-related code, check whether a migration file is required under `backend/migrations/`.
- When editing deployment behavior, update `docs/DEPLOYMENT.md` and related runtime files together.
- Assume the git worktree may already be dirty. Do not revert unrelated user changes.

## Testing And Validation

Backend:

- Prefer `go build ./cmd/api ./cmd/indexer` as a quick validation step.
- Use `go test ./...` when the environment supports executing Go test binaries.
- If Go build cache permissions fail on Windows, prefer setting `GOCACHE` inside the workspace.

Contracts:

- Use Foundry commands from the repo root.
- Typical checks:
  - `forge build`
  - `forge test -vv`

Frontend:

- It is a static app. Prefer validating by serving `frontend/` locally rather than introducing a new framework.

## Deployment Expectations

- Recommended deployment split is:
  - static frontend
  - `api`
  - `indexer`
  - `mysql`
- For container deployment, prefer the existing `docker-compose.yml`.
- Chain runtime configuration must come from the config JSON mounted into the backend container.
- Before claiming deployment is complete, verify `GET /healthz`.

## File-Specific Guidance

- `backend/config/chain.testnet.example.json` is an example config, not production truth.
- `backend/README.md` and `docs/DEPLOYMENT.md` should stay aligned with actual commands and runtime assumptions.
- `frontend/app.js` and `frontend/index.html` may contain user edits; treat them as user-owned unless explicitly asked to modify them.

## Environment Notes

- The workspace is often used on Windows PowerShell.
- Be careful with file encoding. Prefer UTF-8 text files without introducing encoding changes.
- If temporary numbered `.go.*` files appear after tooling issues, treat them as accidental artifacts, not source files.
