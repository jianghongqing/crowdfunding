# Frontend Quick Start

这是一个最小静态页面，用于直连钱包测试合约的四个写方法（Sepolia）：

- `createCampaign`
- `fund`
- `withdraw`
- `refund`

## 使用步骤

1. 先把 `CrowdFund` 合约部署到 Sepolia。
2. 打开 `frontend/index.html`（推荐用本地静态服务）。
3. 页面输入合约地址，点击“连接钱包”。
4. 按需操作创建、捐款、提款、退款。

## Sepolia 准备

- 钱包网络：Sepolia（chainId: `11155111`）
- 测试币：从 Sepolia Faucet 领取
- 区块浏览器：[https://sepolia.etherscan.io](https://sepolia.etherscan.io)

## 本地启动静态服务（可选）

在项目根目录执行：

```bash
npx serve frontend
```

然后打开输出的本地 URL。
