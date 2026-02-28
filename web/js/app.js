/**
 * 主应用 - 可扩展：在此添加新页面的渲染逻辑
 */
document.addEventListener('DOMContentLoaded', () => {
  // 注册页面
  Router.register('/dashboard', renderDashboard);
  Router.register('/stats', renderStats);
  Router.register('/about', renderAbout);

  Router.init();
  checkServerStatus();
  setInterval(checkServerStatus, 10000);
});

// 服务状态检测
async function checkServerStatus() {
  const dot = document.getElementById('serverStatus');
  const text = document.getElementById('serverStatusText');
  try {
    await API.checkHealth();
    dot.className = 'status-dot online';
    text.textContent = '服务正常';
  } catch {
    dot.className = 'status-dot offline';
    text.textContent = '连接失败';
  }
}

// 仪表盘页面
function renderDashboard(container) {
  container.innerHTML = `
    <div class="page-header">
      <h2>仪表盘</h2>
      <p>实时概览手表数据接收情况</p>
    </div>
    <div id="dashboard-stats" class="stats-grid">
      <div class="loading">加载中...</div>
    </div>
  `;
  loadDashboardStats();
}

async function loadDashboardStats() {
  const el = document.getElementById('dashboard-stats');
  if (!el) return;
  try {
    const data = await API.getStats();
    el.innerHTML = `
      <div class="stat-card">
        <div class="stat-label">累计接收数据</div>
        <div class="stat-value">${data.total_received?.toLocaleString() ?? 0}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">接收批次数</div>
        <div class="stat-value">${data.total_batches ?? 0}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">运行时长</div>
        <div class="stat-value" style="font-size:1.25rem">${data.uptime ?? '-'}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">平均速率 (条/秒)</div>
        <div class="stat-value">${(data.avg_rate ?? 0).toFixed(2)}</div>
      </div>
    `;
  } catch (e) {
    el.innerHTML = `<div class="error">加载失败: ${e.message}</div>`;
  }
}

// 统计信息页面（可扩展为更详细的图表等）
function renderStats(container) {
  container.innerHTML = `
    <div class="page-header">
      <h2>统计信息</h2>
      <p>数据接收统计详情</p>
    </div>
    <div class="card">
      <div class="card-title">API 统计</div>
      <div id="stats-detail" class="loading">加载中...</div>
    </div>
    <button class="btn btn-primary" onclick="refreshStats()">🔄 刷新</button>
  `;
  loadStatsDetail();
}

async function loadStatsDetail() {
  const el = document.getElementById('stats-detail');
  if (!el) return;
  try {
    const data = await API.getStats();
    el.innerHTML = `
      <pre style="margin:0;font-size:0.9rem;color:var(--text-muted)">${JSON.stringify(data, null, 2)}</pre>
    `;
  } catch (e) {
    el.innerHTML = `<div class="error">加载失败: ${e.message}</div>`;
  }
}

function refreshStats() {
  const container = document.getElementById('page-content');
  if (container) renderStats(container);
}

// 关于页面
function renderAbout(container) {
  container.innerHTML = `
    <div class="page-header">
      <h2>关于</h2>
      <p>IoT 手表数据管理平台</p>
    </div>
    <div class="about-content card">
      <p>智能手表数据采集与分析后端服务，基于 Gin 框架开发。</p>
      <h3 style="margin:1rem 0 0.5rem">主要接口</h3>
      <ul>
        <li><code>POST /api/v1/data/batch</code> - 批量上传设备数据</li>
        <li><code>GET /api/v1/stats</code> - 获取统计信息</li>
        <li><code>GET /health</code> - 健康检查</li>
      </ul>
      <p style="margin-top:1rem;color:var(--text-muted);font-size:0.9rem">
        访问方式：同一局域网内的设备可通过 <strong>http://&lt;服务器IP&gt;:8080</strong> 访问本管理平台。
      </p>
    </div>
  `;
}
