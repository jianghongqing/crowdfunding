# 生产环境部署指南（云服务器 + 测试链）

本文档以部署到一台云服务器（如阿里云 ECS / 腾讯云 CVM / AWS EC2）为例，使用 Sepolia 测试链作为目标链。

---

## 前置准备

### 服务器要求

| 项目 | 最低配置 | 推荐配置 |
|------|---------|---------|
| CPU | 2 核 | 4 核 |
| 内存 | 4 GB | 8 GB |
| 磁盘 | 40 GB SSD | 80 GB SSD |
| 系统 | Ubuntu 22.04 LTS | Ubuntu 22.04 LTS |
| 网络 | 公网 IP + 开放 80/443 端口 | 同左 |

### 域名与证书（可选但推荐）

- 一个域名（如 `crowdfund.example.com`）
- SSL 证书（Let's Encrypt 免费获取）

### 获取 Sepolia RPC

1. 注册 [Infura](https://infura.io) 或 [Alchemy](https://alchemy.com)（免费计划足够）
2. 创建项目，选择 Sepolia 网络
3. 复制 HTTP 和 WebSocket 端点：
   - HTTP: `https://sepolia.infura.io/v3/YOUR_PROJECT_ID`
   - WS: `wss://sepolia.infura.io/ws/v3/YOUR_PROJECT_ID`

### 获取 Sepolia 测试 ETH

1. 访问 [Sepolia Faucet](https://sepoliafaucet.com) 或 [Alchemy Faucet](https://sepoliafaucet.com)
2. 输入你的钱包地址，领取测试 ETH（用于部署合约和支付 Gas）

---

## 架构总览

```
用户浏览器
    │
    ▼
┌─────────┐     ┌──────────────────┐
│  Nginx  │────▶│  frontend (静态)  │
│ :80/443 │     └──────────────────┘
│         │     ┌──────────────────┐
│         │────▶│   api (:8080)    │
│ /api/*  │     └──────────────────┘
└─────────┘              │
                         ▼
              ┌──────────────────┐
              │  MySQL (:3306)   │◀──── indexer (后台)
              └──────────────────┘
                         │
                    内网访问，不对外暴露
```

---

## 第一步：服务器基础环境

SSH 登录服务器后执行：

```bash
sudo apt update && sudo apt upgrade -y

# 安装 Docker 和 Docker Compose
sudo apt install -y docker.io docker-compose-plugin
sudo systemctl enable docker
sudo systemctl start docker

# 把当前用户加入 docker 组（免 sudo）
sudo usermod -aG docker $USER
# 重新登录生效
exit
```

重新 SSH 登录后验证：

```bash
docker --version
docker compose version
```

### 安装 Foundry（用于部署合约）

```bash
curl -L https://foundry.paradigm.xyz | bash
source ~/.bashrc
foundryup
```

验证：

```bash
forge --version
```

---

## 第二步：拉取项目代码

```bash
cd ~
git clone https://github.com/YOUR_ORG/crowdfunding.git
cd crowdfunding
```

---

## 第三步：部署合约到 Sepolia

**执行目录：** `~/crowdfunding`（仓库根目录）

```bash
forge build
```

部署合约：

```bash
forge script script/DeployCrowdFund.s.sol:DeployCrowdFundScript \
  --rpc-url https://sepolia.infura.io/v3/YOUR_PROJECT_ID \
  --private-key YOUR_DEPLOYER_PRIVATE_KEY \
  --broadcast \
  --verify
```

> **安全警告：**
> - `YOUR_DEPLOYER_PRIVATE_KEY` 是部署者钱包的私钥
> - 该钱包需要有足够的 Sepolia ETH 支付 Gas
> - 部署完成后立即清除终端历史：`history -c`
> - 生产环境建议使用 `--ledger` 或 `cast wallet` 管理密钥

部署成功后记录：

```
合约地址: 0x1234...abcd
部署区块: 7654321
```

可通过 [Sepolia Etherscan](https://sepolia.etherscan.io) 验证合约。

---

## 第四步：准备配置文件

### 4.1 创建链配置 JSON

**执行目录：** `~/crowdfunding`

```bash
cp backend/config/examples/sepolia.example.json backend/config/chain.sepolia.json
```

编辑 `backend/config/chain.sepolia.json`：

```json
{
  "chainName": "sepolia",
  "chainId": 11155111,
  "rpcHttpUrl": "https://sepolia.infura.io/v3/YOUR_PROJECT_ID",
  "rpcWsUrl": "wss://sepolia.infura.io/ws/v3/YOUR_PROJECT_ID",
  "contractAddress": "0x1234...第三步部署得到的合约地址",
  "deploymentStartBlock": 7654321,
  "confirmations": 5
}
```

**字段填写说明：**

| 字段 | 值来源 |
|------|--------|
| `rpcHttpUrl` | Infura/Alchemy 控制台复制的 HTTP 端点 |
| `rpcWsUrl` | Infura/Alchemy 控制台复制的 WS 端点（可选，当前后端未使用） |
| `contractAddress` | 第三步 `forge script` 输出的合约地址 |
| `deploymentStartBlock` | 第三步 `forge script` 输出的区块号 |
| `confirmations` | Sepolia 建议 `5`；代表 indexer 扫到 `最新块 - 5` 时才入库 |

### 4.2 创建 `.env` 文件

```bash
cp .env.example .env
```

编辑 `.env`：

```env
# MySQL
MYSQL_ROOT_PASSWORD=一个强密码_至少16位
MYSQL_DATABASE=crowdfunding

# API 服务
API_PORT=127.0.0.1:8080
CORS_ALLOWED_ORIGINS=https://crowdfund.example.com

# 链配置指向真实文件
CHAIN_CONFIG_HOST_PATH=./backend/config/chain.sepolia.json

# 数据库连接（容器内部网络，host 是 mysql 服务名）
DATABASE_URL=root:一个强密码_至少16位@tcp(mysql:3306)/crowdfunding?charset=utf8mb4&parseTime=true

# 连接池
DB_MAX_OPEN_CONNS=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m
DB_CONN_MAX_IDLE_TIME=5m
```

> **注意：** `MYSQL_ROOT_PASSWORD` 和 `DATABASE_URL` 中的密码必须一致。

---

## 第五步：Docker Compose 启动服务

**执行目录：** `~/crowdfunding`（仓库根目录，`docker-compose.yml` 所在目录）

```bash
# 构建并启动所有服务（后台运行）
docker compose up -d --build
```

查看状态：

```bash
docker compose ps
```

预期输出：

```
NAME                   STATUS
crowdfunding-mysql     Up (healthy)
crowdfunding-api       Up
crowdfunding-indexer   Up
```

### 验证服务

```bash
# API 健康检查
curl http://127.0.0.1:8080/healthz
# 预期: {"status":"ok"}

# 查看 indexer 日志
docker compose logs -f indexer

# 查看 API 日志
docker compose logs -f api
```

---

## 第六步：部署前端

### 方案 A：Nginx 本机托管（推荐）

安装 Nginx：

```bash
sudo apt install -y nginx
```

创建站点配置：

```bash
sudo tee /etc/nginx/sites-available/crowdfunding > /dev/null << 'EOF'
server {
    listen 80;
    server_name crowdfund.example.com;  # 替换成你的域名或服务器 IP

    # 前端静态文件
    root /home/ubuntu/crowdfunding/frontend;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /healthz {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /config {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /campaigns {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
EOF
```

启用站点：

```bash
sudo ln -sf /etc/nginx/sites-available/crowdfunding /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

### 方案 B：使用 HTTPS（Let's Encrypt）

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d crowdfund.example.com
```

Certbot 会自动修改 Nginx 配置并设置自动续期。

### 前端 API 地址

前端默认使用同源 API 路径，不需要修改 `frontend/app.js`。确保 Nginx 已将 `/config`、`/campaigns`、`/healthz` 和 `/api/` 代理到 API 服务。

---

## 第七步：防火墙与安全

### 只开放必要端口

```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

> MySQL (3306) 和 API (8080) **不对外暴露**，仅通过 Nginx 反代和 Docker 内部网络访问。

### 安全建议

- 禁用 MySQL 远程访问（当前 `docker-compose.yml` 默认不暴露 MySQL 端口）
- 使用非 root 用户运行服务
- 定期更新系统和 Docker 镜像
- 不要将 `.env` 和 `chain.sepolia.json` 提交到 Git

---

## 第八步：日常运维

### 查看日志

```bash
docker compose logs -f api       # API 日志
docker compose logs -f indexer   # indexer 日志
docker compose logs -f mysql     # MySQL 日志
```

### 重启服务

```bash
docker compose restart api
docker compose restart indexer
```

### 更新代码和重新部署

```bash
cd ~/crowdfunding
git pull origin main

# 重新构建并重启
docker compose up -d --build
```

### 数据库备份

```bash
# 手动备份
docker exec crowdfunding-mysql mysqldump -u root -p'你的密码' crowdfunding > backup_$(date +%Y%m%d).sql

# 定时备份（添加到 crontab）
crontab -e
# 添加以下行（每天凌晨 3 点备份）：
# 0 3 * * * docker exec crowdfunding-mysql mysqldump -u root -p'你的密码' crowdfunding > /home/ubuntu/backups/crowdfunding_$(date +\%Y\%m\%d).sql
```

### 监控要点

| 监控项 | 方法 | 预期 |
|--------|------|------|
| API 可用性 | `curl /healthz` | 返回 200 |
| Indexer 同步进度 | 查看 indexer 日志 | 区块号持续推进 |
| MySQL 连接数 | `SHOW PROCESSLIST` | 不超过 max_connections |
| 磁盘空间 | `df -h` | 使用率 < 80% |
| RPC 连通性 | `curl $RPC_URL -X POST ...` | 正常返回 |

---

## 完整文件清单

部署后服务器上的关键文件：

```
~/crowdfunding/
├── .env                                    ← 环境变量（不提交 Git）
├── docker-compose.yml                      ← 容器编排
├── backend/
│   ├── config/
│   │   └── chain.sepolia.json             ← 真实链配置（不提交 Git）
│   ├── migrations/                        ← 数据库迁移（自动执行）
│   └── Dockerfile                         ← 后端构建
├── frontend/
│   ├── index.html
│   └── app.js                             ← 默认使用同源 API
└── /etc/nginx/sites-available/crowdfunding ← Nginx 配置
```

---

## 故障排查

### API 启动失败

```bash
docker compose logs api
```

常见原因：
- `DATABASE_URL` 密码与 `MYSQL_ROOT_PASSWORD` 不一致
- MySQL 容器未就绪（等待 healthcheck 通过）
- `CHAIN_CONFIG_HOST_PATH` 指向的文件不存在

### Indexer 无法连接 RPC

```bash
docker compose logs indexer
```

常见原因：
- Infura/Alchemy 限额用尽（换成付费计划或换服务商）
- 服务器无法访问外网（检查安全组出站规则）
- `rpcHttpUrl` 格式错误

### 前端无法调用 API

- 检查 Nginx 是否运行：`sudo systemctl status nginx`
- 检查 API 容器是否运行：`docker compose ps`
- 检查 Nginx 是否代理了 `/config`、`/campaigns`、`/healthz` 和 `/api/`
- 浏览器 F12 查看 Console/Network 错误

### 交易成功但数据不更新

- 查看 indexer 日志确认同步进度
- 确认 `deploymentStartBlock` ≤ 合约部署区块
- 确认 `contractAddress` 与实际部署一致
- 等待 `confirmations` 个块后再检查

---

## 附录：一键部署脚本参考

以下脚本仅供参考，实际使用前请根据环境修改：

```bash
#!/bin/bash
set -e

echo "=== 1. 安装依赖 ==="
sudo apt update && sudo apt install -y docker.io docker-compose-plugin nginx

echo "=== 2. 拉取代码 ==="
cd ~
git clone https://github.com/YOUR_ORG/crowdfunding.git
cd crowdfunding

echo "=== 3. 准备配置 ==="
cp .env.example .env
echo "请编辑 .env 和链配置文件后，按回车继续..."
read

echo "=== 4. 启动服务 ==="
docker compose up -d --build

echo "=== 5. 等待服务就绪 ==="
sleep 15
curl -f http://127.0.0.1:8080/healthz && echo " API OK" || echo " API FAILED"

echo "=== 6. 配置 Nginx ==="
echo "请手动配置 Nginx 站点，参考文档第六步"

echo "=== 部署完成 ==="
```

---

## 相关文档

- [docs/START_FROM_ZERO.md](START_FROM_ZERO.md) — 本地开发启动指南
- [docs/DEPLOYMENT.md](DEPLOYMENT.md) — Docker Compose 部署说明
- [docs/ARCHITECTURE.md](ARCHITECTURE.md) — 系统架构
- [backend/config/README.md](../backend/config/README.md) — 配置文件详细说明
