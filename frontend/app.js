const API_BASE = 'http://localhost:8080';

const ABI = [
  "function getCampaign(uint256 campaignId) view returns (tuple(uint256 id, address creator, string title, uint256 goal, uint256 pledged, uint256 deadline, bool withdrawn))",
  "function campaigns(uint256) view returns (uint256 id, address creator, string title, uint256 goal, uint256 pledged, uint256 deadline, bool withdrawn)",
  "function contributions(uint256, address) view returns (uint256)",
  "function nextCampaignId() view returns (uint256)",
  "function createCampaign(string title, uint256 goal, uint256 duration) returns (uint256)",
  "function fund(uint256 campaignId) payable",
  "function withdraw(uint256 campaignId)",
  "function refund(uint256 campaignId)",
  "event CampaignCreated(uint256 indexed campaignId, address indexed creator, string title, uint256 goal, uint256 deadline)",
  "event Funded(uint256 indexed campaignId, address indexed funder, uint256 amount)",
  "event Withdrawn(uint256 indexed campaignId, address indexed creator, uint256 amount)",
  "event Refunded(uint256 indexed campaignId, address indexed funder, uint256 amount)"
];

const NETWORKS = {
  sepolia: {
    label: "Sepolia 测试网",
    chainIdHex: "0xaa36a7",
    chainName: "Sepolia",
    rpcUrls: ["https://rpc.sepolia.org"],
    blockExplorerUrls: ["https://sepolia.etherscan.io"]
  },
  anvil: {
    label: "本地 Anvil",
    chainIdHex: "0x7a69",
    chainName: "Anvil Local",
    rpcUrls: ["http://127.0.0.1:8545"],
    blockExplorerUrls: []
  }
};

let provider;
let signer;
let contract;

const $ = (id) => document.getElementById(id);

// ---- utilities ----

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function formatEth(wei) {
  return parseFloat(ethers.formatEther(wei)).toFixed(4) + ' ETH';
}

function formatWeiStr(weiStr) {
  if (!weiStr || weiStr === '0') return '0 ETH';
  return formatEth(BigInt(weiStr));
}

function formatDeadline(ts) {
  const d = new Date(Number(ts) * 1000);
  const diff = d - Date.now();
  if (diff <= 0) return d.toLocaleString() + ' (已截止)';
  const days = Math.floor(diff / 86400000);
  const hours = Math.floor((diff % 86400000) / 3600000);
  return d.toLocaleString() + ' (剩余 ' + days + 'd ' + hours + 'h)';
}

function statusLabel(s) {
  const map = {
    active: '进行中',
    goal_reached_pending_withdraw: '已达标(待提款)',
    succeeded_withdrawn: '已提款(成功)',
    failed_refundable: '已失败(可退款)'
  };
  return map[s] || s;
}

function statusBadgeClass(s) {
  const map = {
    active: 'badge-active',
    goal_reached_pending_withdraw: 'badge-goal',
    succeeded_withdrawn: 'badge-success',
    failed_refundable: 'badge-failed'
  };
  return map[s] || '';
}

function log(msg) {
  const el = $("log");
  const now = new Date().toLocaleTimeString();
  el.textContent = '[' + now + '] ' + msg + '\n' + el.textContent;
}

function getNetwork() {
  const key = $("networkSelect").value;
  return NETWORKS[key] || NETWORKS.sepolia;
}

function getAddress() {
  const address = $("contractAddress").value.trim();
  if (!ethers.isAddress(address)) throw new Error("合约地址不合法");
  return address;
}

// ---- API helpers ----

async function apiGet(path) {
  const res = await fetch(API_BASE + path);
  if (!res.ok) {
    const body = await res.json().catch(function() { return {}; });
    throw new Error(body.error || 'HTTP ' + res.status);
  }
  return res.json();
}

// ---- config auto-load ----

async function loadConfig() {
  try {
    const cfg = await apiGet('/config');
    $("contractAddress").value = cfg.contractAddress;
    log('已加载链配置: ' + cfg.chainName + ' (chainId: ' + cfg.chainId + ')');
    return cfg;
  } catch (e) {
    log('加载链配置失败: ' + e.message + ' — 请手动填写合约地址');
    return null;
  }
}

// ---- campaign list ----

async function loadCampaigns() {
  try {
    const data = await apiGet('/campaigns?limit=50');
    renderCampaignList(data);
  } catch (e) {
    $("campaignList").innerHTML = '<div class="muted">加载失败: ' + escapeHtml(e.message) + '</div>';
  }
}

function renderCampaignList(items) {
  var el = $("campaignList");
  if (!items || items.length === 0) {
    el.innerHTML = '<div class="muted">暂无活动</div>';
    return;
  }
  var html = '';
  for (var i = 0; i < items.length; i++) {
    var c = items[i];
    var goal = c.goalWei === '0' ? BigInt(0) : BigInt(c.goalWei);
    var pledged = BigInt(c.pledgedWei);
    var progress = goal === BigInt(0) ? 0 : Number(pledged * BigInt(10000) / goal) / 100;
    if (progress > 100) progress = 100;
    var barColor = c.status === 'failed_refundable' ? '#c62828' :
                   c.status === 'succeeded_withdrawn' ? '#388e3c' : '#4caf50';
    html +=
      '<div class="campaign-row" data-id="' + c.campaignId + '">' +
        '<strong>#' + c.campaignId + '</strong> ' + escapeHtml(c.title) +
        '<span style="float:right;"><span class="badge ' + statusBadgeClass(c.status) + '">' + statusLabel(c.status) + '</span></span>' +
        '<div style="margin-top:4px; font-size:13px; color:#555;">' +
          '进度: ' + formatWeiStr(c.pledgedWei) + ' / ' + formatWeiStr(c.goalWei) +
          ' (' + progress.toFixed(1) + '%)' +
        '</div>' +
        '<div class="progress-bar">' +
          '<div class="progress-fill" style="width:' + progress + '%; background:' + barColor + ';"></div>' +
        '</div>' +
        '<div class="muted" style="margin-top:2px; font-size:11px;">截止: ' + formatDeadline(c.deadline) + '</div>' +
      '</div>';
  }
  el.innerHTML = html;
  el.querySelectorAll('.campaign-row').forEach(function(row) {
    row.addEventListener('click', function() { showCampaignDetail(row.dataset.id); });
  });
}

// ---- campaign detail ----

async function showCampaignDetail(id) {
  try {
    var results = await Promise.all([
      apiGet('/campaigns/' + id),
      apiGet('/campaigns/' + id + '/contributions?limit=100')
    ]);
    var campaign = results[0];
    var contributions = results[1];

    $("campaignDetail").innerHTML =
      '<p><strong>ID:</strong> ' + campaign.campaignId + '</p>' +
      '<p><strong>标题:</strong> ' + escapeHtml(campaign.title) + '</p>' +
      '<p><strong>创建者:</strong> <span style="font-family:monospace; font-size:12px;">' + campaign.creator + '</span></p>' +
      '<p><strong>目标金额:</strong> ' + formatWeiStr(campaign.goalWei) + '</p>' +
      '<p><strong>已筹金额:</strong> ' + formatWeiStr(campaign.pledgedWei) + '</p>' +
      '<p><strong>截止时间:</strong> ' + formatDeadline(campaign.deadline) + '</p>' +
      '<p><strong>状态:</strong> <span class="badge ' + statusBadgeClass(campaign.status) + '">' + statusLabel(campaign.status) + '</span></p>' +
      '<p><strong>创建区块:</strong> ' + (campaign.createdBlock || 'N/A') + '</p>';

    renderContributions(contributions);
    $("campaignListCard").style.display = 'none';
    $("campaignDetailCard").style.display = '';
  } catch (e) {
    log('加载详情失败: ' + e.message);
  }
}

function renderContributions(items) {
  var el = $("contributionList");
  if (!items || items.length === 0) {
    el.innerHTML = '<div class="muted">暂无捐款记录</div>';
    return;
  }
  var html = '';
  for (var i = 0; i < items.length; i++) {
    var c = items[i];
    html +=
      '<div style="border-bottom:1px solid #eee; padding:6px 0; font-size:13px;">' +
        '<span style="font-family:monospace; font-size:11px;">' + c.funder + '</span>' +
        '<span style="float:right; font-weight:bold;">' + formatWeiStr(c.amountWei) + '</span>' +
        '<div class="muted" style="font-size:11px;">tx: ' + (c.txHash ? c.txHash.substring(0, 18) + '...' : 'N/A') + '</div>' +
      '</div>';
  }
  el.innerHTML = html;
}

function backToList() {
  $("campaignListCard").style.display = '';
  $("campaignDetailCard").style.display = 'none';
}

// ---- event listening ----

function setupEventListeners() {
  if (!contract) return;
  contract.off();

  contract.on("CampaignCreated", function(campaignId, creator, title, goal, deadline, eventLog) {
    log('[事件] 新活动创建: #' + campaignId + ' - ' + title + ' (目标: ' + formatEth(goal) + ')');
    loadCampaigns();
  });

  contract.on("Funded", function(campaignId, funder, amount, eventLog) {
    log('[事件] #' + campaignId + ' 收到捐款: ' + formatEth(amount) + ' 来自 ' + funder);
  });

  contract.on("Withdrawn", function(campaignId, creator, amount, eventLog) {
    log('[事件] #' + campaignId + ' 已提款: ' + formatEth(amount));
  });

  contract.on("Refunded", function(campaignId, funder, amount, eventLog) {
    log('[事件] #' + campaignId + ' 退款: ' + formatEth(amount) + ' 给 ' + funder);
  });

  log("事件监听已启动");
}

// ---- wallet ----

async function ensureNetwork() {
  const net = getNetwork();
  const chainId = await window.ethereum.request({ method: "eth_chainId" });
  if (chainId.toLowerCase() === net.chainIdHex.toLowerCase()) return;

  try {
    await window.ethereum.request({
      method: "wallet_switchEthereumChain",
      params: [{ chainId: net.chainIdHex }]
    });
    log('已切换到 ' + net.chainName);
  } catch (err) {
    if (err && err.code === 4902) {
      const params = {
        chainId: net.chainIdHex,
        chainName: net.chainName,
        nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
        rpcUrls: net.rpcUrls
      };
      if (net.blockExplorerUrls.length > 0) {
        params.blockExplorerUrls = net.blockExplorerUrls;
      }
      await window.ethereum.request({
        method: "wallet_addEthereumChain",
        params: [params]
      });
      log('已添加并切换到 ' + net.chainName);
    } else {
      throw new Error('请先手动切换钱包网络到 ' + net.chainName);
    }
  }
}

async function connectWallet() {
  if (!window.ethereum) throw new Error("未检测到钱包，请先安装 MetaMask");
  await window.ethereum.request({ method: "eth_requestAccounts" });
  await ensureNetwork();
  provider = new ethers.BrowserProvider(window.ethereum);
  signer = await provider.getSigner();
  const address = getAddress();
  contract = new ethers.Contract(address, ABI, signer);
  setupEventListeners();
  const network = await provider.getNetwork();
  const net = getNetwork();
  $("walletInfo").textContent = '已连接: ' + (await signer.getAddress()) + ' | ' + net.label + ' | chainId: ' + network.chainId.toString();
  log("钱包连接成功");
}

function ensureContract() {
  if (!contract) throw new Error("请先输入合约地址并连接钱包");
}

// ---- write operations ----

async function createCampaign() {
  ensureContract();
  const title = $("createTitle").value.trim();
  const goalEth = $("createGoalEth").value.trim();
  const duration = $("createDuration").value.trim();

  if (!title) throw new Error("标题不能为空");
  if (!goalEth) throw new Error("目标金额不能为空");
  if (!duration) throw new Error("持续时间不能为空");

  const tx = await contract.createCampaign(title, ethers.parseEther(goalEth), BigInt(duration));
  log('create 交易已发送: ' + tx.hash);
  const receipt = await tx.wait();
  log('create 成功, block: ' + receipt.blockNumber + ' — 刷新列表查看');
}

async function fundCampaign() {
  ensureContract();
  const id = $("fundCampaignId").value.trim();
  const amountEth = $("fundAmountEth").value.trim();
  if (!id || !amountEth) throw new Error("campaignId 和金额不能为空");

  const tx = await contract.fund(BigInt(id), { value: ethers.parseEther(amountEth) });
  log('fund 交易已发送: ' + tx.hash);
  const receipt = await tx.wait();
  log('fund 成功, block: ' + receipt.blockNumber);
}

async function withdrawCampaign() {
  ensureContract();
  const id = $("withdrawCampaignId").value.trim();
  if (!id) throw new Error("campaignId 不能为空");

  const tx = await contract.withdraw(BigInt(id));
  log('withdraw 交易已发送: ' + tx.hash);
  const receipt = await tx.wait();
  log('withdraw 成功, block: ' + receipt.blockNumber);
}

async function refundCampaign() {
  ensureContract();
  const id = $("refundCampaignId").value.trim();
  if (!id) throw new Error("campaignId 不能为空");

  const tx = await contract.refund(BigInt(id));
  log('refund 交易已发送: ' + tx.hash);
  const receipt = await tx.wait();
  log('refund 成功, block: ' + receipt.blockNumber);
}

// ---- event binding ----

function bind(id, fn) {
  $(id).addEventListener("click", async function() {
    try {
      await fn();
    } catch (err) {
      log('错误: ' + (err?.reason || err?.shortMessage || err?.message || String(err)));
    }
  });
}

bind("connectBtn", connectWallet);
bind("createBtn", createCampaign);
bind("fundBtn", fundCampaign);
bind("withdrawBtn", withdrawCampaign);
bind("refundBtn", refundCampaign);
bind("loadCampaignsBtn", loadCampaigns);
bind("backToListBtn", backToList);

$("networkSelect").addEventListener("change", function() {
  contract = null;
  const net = getNetwork();
  log('已选择网络: ' + net.label + '，请重新点击连接钱包');
  $("walletInfo").textContent = "网络已更改，请重新点击连接钱包";
});

if (window.ethereum) {
  window.ethereum.on("chainChanged", function(chainId) {
    const net = getNetwork();
    if (chainId.toLowerCase() !== net.chainIdHex.toLowerCase()) {
      log('当前链与所选「' + net.label + '」不一致，交易可能失败');
    } else {
      log('当前网络: ' + net.chainName);
    }
    contract = null;
    $("walletInfo").textContent = "网络已切换，请重新点击连接钱包";
  });
}

// ---- startup ----

(async function init() {
  try {
    await loadConfig();
    await loadCampaigns();
  } catch (e) {
    log('初始化失败: ' + e.message);
  }
})();
