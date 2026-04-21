const ABI = [
  "function createCampaign(string title, uint256 goal, uint256 duration) returns (uint256)",
  "function fund(uint256 campaignId) payable",
  "function withdraw(uint256 campaignId)",
  "function refund(uint256 campaignId)"
];
const SEPOLIA_CHAIN_ID_HEX = "0xaa36a7";
const SEPOLIA_CHAIN_ID_DEC = 11155111n;

let provider;
let signer;
let contract;

const $ = (id) => document.getElementById(id);

function log(msg) {
  const el = $("log");
  const now = new Date().toLocaleTimeString();
  el.textContent = `[${now}] ${msg}\n` + el.textContent;
}

function getAddress() {
  const address = $("contractAddress").value.trim();
  if (!ethers.isAddress(address)) throw new Error("合约地址不合法");
  return address;
}

async function connectWallet() {
  if (!window.ethereum) throw new Error("未检测到钱包，请先安装 MetaMask");
  await window.ethereum.request({ method: "eth_requestAccounts" });
  await ensureSepoliaNetwork();
  provider = new ethers.BrowserProvider(window.ethereum);
  signer = await provider.getSigner();
  const address = getAddress();
  contract = new ethers.Contract(address, ABI, signer);
  const network = await provider.getNetwork();
  $("walletInfo").textContent = `已连接: ${await signer.getAddress()} | chainId: ${network.chainId.toString()}`;
  log("钱包连接成功");
}

async function ensureSepoliaNetwork() {
  const chainId = await window.ethereum.request({ method: "eth_chainId" });
  if (chainId.toLowerCase() === SEPOLIA_CHAIN_ID_HEX) return;

  try {
    await window.ethereum.request({
      method: "wallet_switchEthereumChain",
      params: [{ chainId: SEPOLIA_CHAIN_ID_HEX }]
    });
    log("已切换到 Sepolia");
  } catch (err) {
    if (err && err.code === 4902) {
      await window.ethereum.request({
        method: "wallet_addEthereumChain",
        params: [{
          chainId: SEPOLIA_CHAIN_ID_HEX,
          chainName: "Sepolia",
          nativeCurrency: { name: "SepoliaETH", symbol: "ETH", decimals: 18 },
          rpcUrls: ["https://rpc.sepolia.org"],
          blockExplorerUrls: ["https://sepolia.etherscan.io"]
        }]
      });
      log("已添加并切换到 Sepolia");
    } else {
      throw new Error("请先手动切换钱包网络到 Sepolia");
    }
  }
}

function ensureContract() {
  if (!contract) throw new Error("请先输入合约地址并连接钱包");
}

async function createCampaign() {
  ensureContract();
  const title = $("createTitle").value.trim();
  const goalEth = $("createGoalEth").value.trim();
  const duration = $("createDuration").value.trim();

  if (!title) throw new Error("标题不能为空");
  if (!goalEth) throw new Error("目标金额不能为空");
  if (!duration) throw new Error("持续时间不能为空");

  const tx = await contract.createCampaign(title, ethers.parseEther(goalEth), BigInt(duration));
  log(`create 交易已发送: ${tx.hash}`);
  const receipt = await tx.wait();
  log(`create 成功, block: ${receipt.blockNumber}`);
}

async function fundCampaign() {
  ensureContract();
  const id = $("fundCampaignId").value.trim();
  const amountEth = $("fundAmountEth").value.trim();
  if (!id || !amountEth) throw new Error("campaignId 和金额不能为空");

  const tx = await contract.fund(BigInt(id), { value: ethers.parseEther(amountEth) });
  log(`fund 交易已发送: ${tx.hash}`);
  const receipt = await tx.wait();
  log(`fund 成功, block: ${receipt.blockNumber}`);
}

async function withdrawCampaign() {
  ensureContract();
  const id = $("withdrawCampaignId").value.trim();
  if (!id) throw new Error("campaignId 不能为空");

  const tx = await contract.withdraw(BigInt(id));
  log(`withdraw 交易已发送: ${tx.hash}`);
  const receipt = await tx.wait();
  log(`withdraw 成功, block: ${receipt.blockNumber}`);
}

async function refundCampaign() {
  ensureContract();
  const id = $("refundCampaignId").value.trim();
  if (!id) throw new Error("campaignId 不能为空");

  const tx = await contract.refund(BigInt(id));
  log(`refund 交易已发送: ${tx.hash}`);
  const receipt = await tx.wait();
  log(`refund 成功, block: ${receipt.blockNumber}`);
}

function bind(id, fn) {
  $(id).addEventListener("click", async () => {
    try {
      await fn();
    } catch (err) {
      log(`错误: ${err?.reason || err?.shortMessage || err?.message || String(err)}`);
    }
  });
}

bind("connectBtn", connectWallet);
bind("createBtn", createCampaign);
bind("fundBtn", fundCampaign);
bind("withdrawBtn", withdrawCampaign);
bind("refundBtn", refundCampaign);

if (window.ethereum) {
  window.ethereum.on("chainChanged", (chainId) => {
    if (chainId.toLowerCase() !== SEPOLIA_CHAIN_ID_HEX) {
      log("当前网络不是 Sepolia，交易可能失败");
    } else {
      log("当前网络: Sepolia");
    }
    contract = null;
    $("walletInfo").textContent = "网络已切换，请重新点击连接钱包";
  });
}
