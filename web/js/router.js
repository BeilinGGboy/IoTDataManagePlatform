/**
 * 简单 Hash 路由 - 可扩展：在此注册新页面
 */
const Router = {
  routes: {},

  register(path, renderFn) {
    this.routes[path] = renderFn;
  },

  init() {
    window.addEventListener('hashchange', () => this.handleRoute());
    window.addEventListener('load', () => this.handleRoute());
  },

  handleRoute() {
    const hash = window.location.hash.slice(1) || '/dashboard';
    const path = hash.startsWith('/') ? hash : '/' + hash;
    const [basePath] = path.split('/').filter(Boolean);
    const route = '/' + (basePath || 'dashboard');

    const renderFn = this.routes[route] || this.routes['/dashboard'];
    const container = document.getElementById('page-content');
    if (container && renderFn) {
      container.innerHTML = '';
      renderFn(container);
    }

    // 更新导航激活状态
    document.querySelectorAll('.nav-item').forEach(el => {
      el.classList.toggle('active', el.dataset.page === (basePath || 'dashboard'));
    });
  }
};
