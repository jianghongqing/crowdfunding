#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd)"
BACKEND_DIR="${ROOT_DIR}/backend"

cd "${ROOT_DIR}"

echo "Generating ABI and BIN from Foundry artifact..."
node -e "const fs=require('fs'); const d=JSON.parse(fs.readFileSync('out/CrowdFund.sol/CrowdFund.json','utf8')); fs.writeFileSync('backend/contracts/crowdfund/CrowdFund.abi', JSON.stringify(d.abi)); fs.writeFileSync('backend/contracts/crowdfund/CrowdFund.bin', d.bytecode.object);"

echo "Generating Go binding via abigen..."
ABIGEN_BIN="${ABIGEN_BIN:-${HOME}/go/bin/linux_amd64/abigen}"
"${ABIGEN_BIN}" \
  --abi "${BACKEND_DIR}/contracts/crowdfund/CrowdFund.abi" \
  --bin "${BACKEND_DIR}/contracts/crowdfund/CrowdFund.bin" \
  --pkg crowdfund \
  --type CrowdFund \
  --out "${BACKEND_DIR}/contracts/crowdfund/crowdfund.go"

echo "Done."
