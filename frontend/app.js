/**
 * CrowdFund DApp 前端主逻辑。
 *
 * 架构：纯静态页面 + ethers.js，不依赖构建工具。
 * - 钱包交互：通过 MetaMask 等 window.ethereum 注入器直接签名发交易
 * - 数据展示：调用后端 API（/campaigns, /config 等）获取 indexer 同步的链下数据
 * - 实时事件：通过 ethers 合约事件监听，在交易确认后即时通知用户
 */
'use strict';

// 后端 API 地址，Docker Compose 部署时由 nginx 反向代理，无需修改
const API_BASE = 'http://localhost:8080';

// CrowdFund 合约的人类可读 ABI（ethers v6 格式）
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

// 支持的网络配置，包含链信息和区块浏览器 URL 生成器
const NETWORKS = {
  sepolia: {
    label: "Sepolia 测试网",
    chainIdHex: "0xaa36a7",
    chainName: "Sepolia",
    rpcUrls: ["https://rpc.sepolia.org"],
    blockExplorerUrls: ["https://sepolia.etherscan.io"],
    explorerTx: (hash) => `https://sepolia.etherscan.io/tx/${hash}`,
    explorerAddr: (addr) => `https://sepolia.etherscan.io/address/${addr}`,
  },
  anvil: {
    label: "本地 Anvil",
    chainIdHex: "0x7a69",
    chainName: "Anvil Local",
    rpcUrls: ["http://127.0.0.1:8545"],
    blockExplorerUrls: [],
    explorerTx: () => null,
    explorerAddr: () => null,
  }
};

// ---- 全局状态 ----
// 连接钱包后由 connectWallet() 初始化，断开时重置为 null
let provider = null;   // ethers.BrowserProvider
let signer = null;     // 当前钱包签名器
let contract = null;   // 绑定了 signer 的合约实例，用于发交易
let currentCampaigns = []; // 缓存当前活动列表，供 updateStats 使用

// ---- DOM Helpers ----
const $ = (id) => document.getElementById(id);
const $$ = (sel) => document.querySelectorAll(sel);

// ---- Toast 通知系统 ----
// 替代原始日志面板，4秒自动消失，支持 success/error/info/warning 四种类型
function toast(message, type = 'info') {
  const container = $('toastContainer');
  const colors = {
    success: 'bg-green-600',
    error: 'bg-red-600',
    info: 'bg-brand-600',
    warning: 'bg-orange-500',
  };
  const icons = {
    success: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>',
    error: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>',
    info: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>',
    warning: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z"/>',
  };

  const el = document.createElement('div');
  el.className = `toast-enter flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg text-white text-sm ${colors[type] || colors.info}`;
  el.innerHTML = `<svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">${icons[type] || icons.info}</svg><span class="flex-1">${escapeHtml(message)}</span>`;
  container.appendChild(el);

  setTimeout(() => {
    el.style.opacity = '0';
    el.style.transition = 'opacity 0.3s';
    setTimeout(() => el.remove(), 300);
  }, 4000);
}

// ---- Utilities ----
// 利用 DOM textContent 自动转义，防止 XSS
function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function formatEth(wei) {
  return parseFloat(ethers.formatEther(wei)).toFixed(4);
}

function formatWeiStr(weiStr) {
  if (!weiStr || weiStr === '0') return '0';
  return formatEth(BigInt(weiStr));
}

function shortAddress(addr) {
  if (!addr) return '';
  return addr.slice(0, 6) + '...' + addr.slice(-4);
}

function formatDeadline(ts) {
  const d = new Date(Number(ts) * 1000);
  const diff = d - Date.now();
  if (diff <= 0) return { text: '已截止', expired: true, relative: '' };
  const days = Math.floor(diff / 86400000);
  const hours = Math.floor((diff % 86400000) / 3600000);
  const mins = Math.floor((diff % 3600000) / 60000);
  let relative = '';
  if (days > 0) relative = `${days}天 ${hours}小时`;
  else if (hours > 0) relative = `${hours}小时 ${mins}分钟`;
  else relative = `${mins}分钟`;
  return { text: d.toLocaleString('zh-CN'), expired: false, relative };
}

function statusLabel(s) {
  const map = {
    active: '进行中',
    goal_reached_pending_withdraw: '已达标',
    succeeded_withdrawn: '已成功',
    failed_refundable: '已失败',
  };
  return map[s] || s;
}

// 每种活动状态对应的 Tailwind CSS 类名配置（背景、文字、圆点、进度条颜色）
function statusConfig(s) {
  const configs = {
    active: { bg: 'bg-blue-100', text: 'text-blue-700', dot: 'bg-blue-500', bar: 'bg-blue-500' },
    goal_reached_pending_withdraw: { bg: 'bg-amber-100', text: 'text-amber-700', dot: 'bg-amber-500', bar: 'bg-amber-500' },
    succeeded_withdrawn: { bg: 'bg-green-100', text: 'text-green-700', dot: 'bg-green-500', bar: 'bg-green-500' },
    failed_refundable: { bg: 'bg-red-100', text: 'text-red-700', dot: 'bg-red-500', bar: 'bg-red-400' },
  };
  return configs[s] || configs.active;
}

function getNetwork() {
  return NETWORKS[$('networkSelect').value] || NETWORKS.sepolia;
}

function getContractAddress() {
  const el = $('contractAddress');
  return el ? el.value : '';
}

// ---- API Helpers ----
async function apiGet(path) {
  const res = await fetch(API_BASE + path);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json();
}

// ---- Tab Navigation ----
function switchTab(tabName) {
  $$('.tab-btn').forEach(btn => {
    if (btn.dataset.tab === tabName) {
      btn.classList.add('bg-white', 'text-gray-900', 'shadow-sm');
      btn.classList.remove('text-gray-600');
    } else {
      btn.classList.remove('bg-white', 'text-gray-900', 'shadow-sm');
      btn.classList.add('text-gray-600');
    }
  });

  $$('.tab-panel').forEach(p => p.classList.add('hidden'));
  const panel = $('panel' + tabName.charAt(0).toUpperCase() + tabName.slice(1));
  if (panel) panel.classList.remove('hidden');

  // Hide detail panel when switching tabs
  $('panelDetail').classList.add('hidden');
}

// ---- Config & Stats ----
let contractAddress = '';

async function loadConfig() {
  try {
    const cfg = await apiGet('/config');
    contractAddress = cfg.contractAddress;
    $('contractDisplay').textContent = shortAddress(cfg.contractAddress);
    toast(`已连接: ${cfg.chainName} (chainId: ${cfg.chainId})`, 'info');
    return cfg;
  } catch (e) {
    toast('配置加载失败，请确认后端已启动', 'warning');
    return null;
  }
}

function updateStats(campaigns) {
  const total = campaigns.length;
  const active = campaigns.filter(c => c.status === 'active' || c.status === 'goal_reached_pending_withdraw').length;
  const success = campaigns.filter(c => c.status === 'succeeded_withdrawn').length;
  $('statTotal').textContent = total;
  $('statActive').textContent = active;
  $('statSuccess').textContent = success;
}

// ---- Campaign List ----
async function loadCampaigns() {
  try {
    const resp = await apiGet('/campaigns?limit=50');
    const items = Array.isArray(resp) ? resp : (resp.data || []);
    currentCampaigns = items;
    renderCampaignList(currentCampaigns);
    updateStats(currentCampaigns);
  } catch (e) {
    $('campaignList').innerHTML = `<div class="col-span-full text-center py-8 text-gray-500">${escapeHtml(e.message)}</div>`;
  }
}

function renderCampaignList(items) {
  const el = $('campaignList');
  const empty = $('emptyState');

  if (!items || items.length === 0) {
    el.innerHTML = '';
    empty.classList.remove('hidden');
    return;
  }

  empty.classList.add('hidden');
  let html = '';

  for (const c of items) {
    const goal = BigInt(c.goalWei || '0');
    const pledged = BigInt(c.pledgedWei || '0');
    const progress = goal === 0n ? 0 : Math.min(Number(pledged * 10000n / goal) / 100, 100);
    const sc = statusConfig(c.status);
    const deadline = formatDeadline(c.deadline);

    html += `
      <div class="card-hover bg-white rounded-xl border border-gray-200 shadow-sm p-5 cursor-pointer animate-fade-in" data-id="${c.campaignId}" onclick="showCampaignDetail(${c.campaignId})">
        <div class="flex items-start justify-between mb-3">
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-gray-900 truncate">${escapeHtml(c.title)}</h3>
            <p class="text-xs text-gray-400 mt-0.5 font-mono">#${c.campaignId}</p>
          </div>
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium ${sc.bg} ${sc.text}">
            <span class="w-1.5 h-1.5 rounded-full ${sc.dot}"></span>
            ${statusLabel(c.status)}
          </span>
        </div>

        <div class="mb-3">
          <div class="flex justify-between text-sm mb-1">
            <span class="text-gray-600">${formatWeiStr(c.pledgedWei)} ETH</span>
            <span class="text-gray-400">目标 ${formatWeiStr(c.goalWei)} ETH</span>
          </div>
          <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
            <div class="h-full rounded-full transition-all duration-500 ${sc.bar}" style="width: ${progress}%"></div>
          </div>
          <div class="text-xs text-gray-400 mt-1 text-right">${progress.toFixed(1)}%</div>
        </div>

        <div class="flex items-center justify-between text-xs text-gray-500">
          <span class="font-mono">${shortAddress(c.creator)}</span>
          <span>${deadline.expired ? '⏱ 已截止' : '⏱ 剩余 ' + deadline.relative}</span>
        </div>
      </div>`;
  }

  el.innerHTML = html;
}

// ---- Campaign Detail ----
async function showCampaignDetail(id) {
  try {
    const [campaignResp, contribResp] = await Promise.all([
      apiGet('/campaigns/' + id),
      apiGet('/campaigns/' + id + '/contributions?limit=100'),
    ]);
    const campaign = campaignResp.data ? campaignResp.data : campaignResp;
    const contributions = Array.isArray(contribResp) ? contribResp : (contribResp.data || []);

    const goal = BigInt(campaign.goalWei || '0');
    const pledged = BigInt(campaign.pledgedWei || '0');
    const progress = goal === 0n ? 0 : Math.min(Number(pledged * 10000n / goal) / 100, 100);
    const sc = statusConfig(campaign.status);
    const deadline = formatDeadline(campaign.deadline);
    const net = getNetwork();

    const explorerLink = net.explorerAddr ? net.explorerAddr(campaign.creator) : null;
    const creatorDisplay = explorerLink
      ? `<a href="${explorerLink}" target="_blank" class="text-brand-600 hover:underline font-mono text-sm">${campaign.creator}</a>`
      : `<span class="font-mono text-sm text-gray-700">${campaign.creator}</span>`;

    $('campaignDetail').innerHTML = `
      <div class="flex items-start justify-between mb-6">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">${escapeHtml(campaign.title)}</h2>
          <p class="text-sm text-gray-500 mt-1">活动 #${campaign.campaignId}</p>
        </div>
        <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium ${sc.bg} ${sc.text}">
          <span class="w-2 h-2 rounded-full ${sc.dot}"></span>
          ${statusLabel(campaign.status)}
        </span>
      </div>

      <div class="mb-6">
        <div class="flex justify-between text-sm mb-2">
          <span class="font-medium text-gray-700">${formatWeiStr(campaign.pledgedWei)} ETH 已筹集</span>
          <span class="text-gray-500">目标 ${formatWeiStr(campaign.goalWei)} ETH</span>
        </div>
        <div class="h-3 bg-gray-100 rounded-full overflow-hidden">
          <div class="h-full rounded-full transition-all duration-700 ${sc.bar}" style="width: ${progress}%"></div>
        </div>
        <div class="text-sm text-gray-500 mt-2 text-right">${progress.toFixed(1)}% 完成</div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 bg-gray-50 rounded-lg">
        <div>
          <div class="text-xs text-gray-500 uppercase tracking-wide">创建者</div>
          <div class="mt-1">${creatorDisplay}</div>
        </div>
        <div>
          <div class="text-xs text-gray-500 uppercase tracking-wide">截止时间</div>
          <div class="mt-1 text-sm text-gray-700">${deadline.text}</div>
          ${!deadline.expired ? `<div class="text-xs text-brand-600 font-medium">剩余 ${deadline.relative}</div>` : '<div class="text-xs text-red-600 font-medium">已截止</div>'}
        </div>
        <div>
          <div class="text-xs text-gray-500 uppercase tracking-wide">创建区块</div>
          <div class="mt-1 text-sm text-gray-700">${campaign.createdBlock || 'N/A'}</div>
        </div>
        <div>
          <div class="text-xs text-gray-500 uppercase tracking-wide">资金状态</div>
          <div class="mt-1 text-sm text-gray-700">${campaign.withdrawn ? '已提取' : '合约托管中'}</div>
        </div>
      </div>

      ${campaign.status === 'active' ? `
        <div class="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <h4 class="font-medium text-blue-900 mb-2">支持此活动</h4>
          <div class="flex gap-2">
            <input id="detailFundAmount" type="number" step="0.001" min="0.001" placeholder="ETH 金额" class="flex-1 px-3 py-2 border border-blue-200 rounded-lg text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
            <button onclick="fundFromDetail(${campaign.campaignId})" class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors">捐款</button>
          </div>
        </div>
      ` : ''}
    `;

    renderContributions(contributions);
    $$('.tab-panel').forEach(p => p.classList.add('hidden'));
    $('panelDetail').classList.remove('hidden');
  } catch (e) {
    toast('加载活动详情失败: ' + e.message, 'error');
  }
}

function renderContributions(items) {
  const el = $('contributionList');
  if (!items || items.length === 0) {
    el.innerHTML = '<div class="p-6 text-center text-sm text-gray-500">暂无捐款记录</div>';
    return;
  }

  const net = getNetwork();
  let html = `
    <table class="w-full">
      <thead class="bg-gray-50 border-b border-gray-200">
        <tr>
          <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">捐款人</th>
          <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">金额</th>
          <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase hidden sm:table-cell">交易</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-100">`;

  for (const c of items) {
    const txLink = net.explorerTx && c.txHash ? net.explorerTx(c.txHash) : null;
    const txDisplay = txLink
      ? `<a href="${txLink}" target="_blank" class="text-brand-600 hover:underline">${c.txHash.slice(0, 10)}...</a>`
      : `<span class="text-gray-400">${c.txHash ? c.txHash.slice(0, 10) + '...' : 'N/A'}</span>`;

    html += `
      <tr class="hover:bg-gray-50">
        <td class="px-4 py-3 text-sm font-mono text-gray-700">${shortAddress(c.funder)}</td>
        <td class="px-4 py-3 text-sm text-right font-medium text-gray-900">${formatWeiStr(c.amountWei)} ETH</td>
        <td class="px-4 py-3 text-sm text-right font-mono hidden sm:table-cell">${txDisplay}</td>
      </tr>`;
  }

  html += '</tbody></table>';
  el.innerHTML = html;
}

function backToList() {
  $('panelDetail').classList.add('hidden');
  $('panelCampaigns').classList.remove('hidden');
  switchTab('campaigns');
}

// ---- Fund from detail page ----
async function fundFromDetail(campaignId) {
  if (!contract) {
    toast('请先连接钱包', 'warning');
    return;
  }
  const amount = $('detailFundAmount').value.trim();
  if (!amount || parseFloat(amount) <= 0) {
    toast('请输入有效金额', 'warning');
    return;
  }
  try {
    const tx = await contract.fund(BigInt(campaignId), { value: ethers.parseEther(amount) });
    toast('交易已提交，等待确认...', 'info');
    await tx.wait();
    toast('捐款成功!', 'success');
    showCampaignDetail(campaignId);
    loadCampaigns();
  } catch (err) {
    toast(err?.reason || err?.shortMessage || err?.message || '交易失败', 'error');
  }
}

// Make it globally accessible for onclick
window.showCampaignDetail = showCampaignDetail;
window.fundFromDetail = fundFromDetail;

// ---- Event Listening ----
function setupEventListeners() {
  if (!contract) return;
  contract.off();

  contract.on("CampaignCreated", (campaignId, creator, title) => {
    toast(`新活动创建: #${campaignId} - ${title}`, 'success');
    loadCampaigns();
  });

  contract.on("Funded", (campaignId, funder, amount) => {
    toast(`#${campaignId} 收到 ${formatEth(amount)} ETH 捐款`, 'info');
  });

  contract.on("Withdrawn", (campaignId, creator, amount) => {
    toast(`#${campaignId} 已提款 ${formatEth(amount)} ETH`, 'success');
    loadCampaigns();
  });

  contract.on("Refunded", (campaignId, funder, amount) => {
    toast(`#${campaignId} 退款 ${formatEth(amount)} ETH`, 'info');
  });
}

// ---- Wallet ----
async function ensureNetwork() {
  const net = getNetwork();
  const chainId = await window.ethereum.request({ method: "eth_chainId" });
  if (chainId.toLowerCase() === net.chainIdHex.toLowerCase()) return;

  try {
    await window.ethereum.request({
      method: "wallet_switchEthereumChain",
      params: [{ chainId: net.chainIdHex }],
    });
  } catch (err) {
    if (err && err.code === 4902) {
      const params = {
        chainId: net.chainIdHex,
        chainName: net.chainName,
        nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
        rpcUrls: net.rpcUrls,
      };
      if (net.blockExplorerUrls.length > 0) params.blockExplorerUrls = net.blockExplorerUrls;
      await window.ethereum.request({ method: "wallet_addEthereumChain", params: [params] });
    } else {
      throw new Error('请先切换钱包网络到 ' + net.chainName);
    }
  }
}

async function connectWallet() {
  if (!window.ethereum) {
    toast('未检测到钱包，请先安装 MetaMask', 'error');
    return;
  }

  try {
    $('connectBtn').disabled = true;
    $('connectBtnText').textContent = '连接中...';

    await window.ethereum.request({ method: "eth_requestAccounts" });
    await ensureNetwork();

    provider = new ethers.BrowserProvider(window.ethereum);
    signer = await provider.getSigner();

    if (!contractAddress) {
      toast('合约地址未加载，请检查后端配置', 'error');
      return;
    }

    contract = new ethers.Contract(contractAddress, ABI, signer);
    setupEventListeners();

    const address = await signer.getAddress();
    const net = getNetwork();

    $('walletBar').classList.remove('hidden');
    $('walletAddress').textContent = shortAddress(address);
    $('walletNetwork').textContent = net.label;
    $('connectBtnText').textContent = shortAddress(address);

    toast('钱包已连接', 'success');
  } catch (err) {
    toast(err?.message || '钱包连接失败', 'error');
    $('connectBtnText').textContent = '连接钱包';
  } finally {
    $('connectBtn').disabled = false;
  }
}

function ensureContract() {
  if (!contract) throw new Error('请先连接钱包');
}

// ---- Write Operations ----
function setLoading(btn, loading) {
  btn.disabled = loading;
  if (loading) btn.classList.add('opacity-50');
  else btn.classList.remove('opacity-50');
}

async function createCampaign(e) {
  e.preventDefault();
  ensureContract();

  const title = $('createTitle').value.trim();
  const goalEth = $('createGoalEth').value.trim();
  const duration = $('createDuration').value;

  if (!title) throw new Error('标题不能为空');
  if (!goalEth || parseFloat(goalEth) <= 0) throw new Error('目标金额无效');

  const btn = $('createBtn');
  setLoading(btn, true);
  $('createBtnText').textContent = '交易确认中...';
  $('createSpinner').classList.remove('hidden');

  try {
    const tx = await contract.createCampaign(title, ethers.parseEther(goalEth), BigInt(duration));
    toast('交易已提交，等待链上确认...', 'info');
    await tx.wait();
    toast('活动创建成功!', 'success');

    $('createTitle').value = '';
    $('createGoalEth').value = '';
    loadCampaigns();
    switchTab('campaigns');
  } finally {
    setLoading(btn, false);
    $('createBtnText').textContent = '发起活动';
    $('createSpinner').classList.add('hidden');
  }
}

async function fundCampaign() {
  ensureContract();
  const id = $('fundCampaignId').value.trim();
  const amountEth = $('fundAmountEth').value.trim();
  if (!id || !amountEth) throw new Error('活动 ID 和金额不能为空');

  const btn = $('fundBtn');
  setLoading(btn, true);
  try {
    const tx = await contract.fund(BigInt(id), { value: ethers.parseEther(amountEth) });
    toast('捐款交易已提交...', 'info');
    await tx.wait();
    toast('捐款成功!', 'success');
    $('fundCampaignId').value = '';
    $('fundAmountEth').value = '';
    loadCampaigns();
  } finally {
    setLoading(btn, false);
  }
}

async function withdrawCampaign() {
  ensureContract();
  const id = $('withdrawCampaignId').value.trim();
  if (!id) throw new Error('活动 ID 不能为空');

  const btn = $('withdrawBtn');
  setLoading(btn, true);
  try {
    const tx = await contract.withdraw(BigInt(id));
    toast('提款交易已提交...', 'info');
    await tx.wait();
    toast('提款成功!', 'success');
    $('withdrawCampaignId').value = '';
    loadCampaigns();
  } finally {
    setLoading(btn, false);
  }
}

async function refundCampaign() {
  ensureContract();
  const id = $('refundCampaignId').value.trim();
  if (!id) throw new Error('活动 ID 不能为空');

  const btn = $('refundBtn');
  setLoading(btn, true);
  try {
    const tx = await contract.refund(BigInt(id));
    toast('退款交易已提交...', 'info');
    await tx.wait();
    toast('退款成功!', 'success');
    $('refundCampaignId').value = '';
    loadCampaigns();
  } finally {
    setLoading(btn, false);
  }
}

// ---- Event Binding ----
function bind(id, fn) {
  const el = $(id);
  if (!el) return;
  el.addEventListener('click', async (e) => {
    try {
      await fn(e);
    } catch (err) {
      toast(err?.reason || err?.shortMessage || err?.message || String(err), 'error');
    }
  });
}

bind('connectBtn', connectWallet);
bind('fundBtn', fundCampaign);
bind('withdrawBtn', withdrawCampaign);
bind('refundBtn', refundCampaign);
bind('refreshBtn', loadCampaigns);
bind('backToListBtn', backToList);

$('createForm').addEventListener('submit', async (e) => {
  try {
    await createCampaign(e);
  } catch (err) {
    toast(err?.reason || err?.shortMessage || err?.message || String(err), 'error');
  }
});

$$('.tab-btn').forEach(btn => {
  btn.addEventListener('click', () => switchTab(btn.dataset.tab));
});

$('networkSelect').addEventListener('change', () => {
  contract = null;
  signer = null;
  $('walletBar').classList.add('hidden');
  $('connectBtnText').textContent = '连接钱包';
  toast('网络已切换，请重新连接钱包', 'warning');
});

if (window.ethereum) {
  window.ethereum.on('chainChanged', () => {
    contract = null;
    signer = null;
    $('walletBar').classList.add('hidden');
    $('connectBtnText').textContent = '连接钱包';
    toast('网络已变更，请重新连接', 'warning');
  });

  window.ethereum.on('accountsChanged', (accounts) => {
    if (accounts.length === 0) {
      contract = null;
      signer = null;
      $('walletBar').classList.add('hidden');
      $('connectBtnText').textContent = '连接钱包';
    } else {
      connectWallet();
    }
  });
}

// ---- Startup ----
(async function init() {
  await loadConfig();
  await loadCampaigns();
})();
