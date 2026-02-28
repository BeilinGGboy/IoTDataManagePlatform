/**
 * API 调用模块 - 可扩展：在此添加新的 API 方法
 */
const API = {
  baseUrl: '', // 同源，使用相对路径

  async getStats() {
    const res = await fetch(`${this.baseUrl}/api/v1/stats`);
    if (!res.ok) throw new Error('获取统计失败');
    return res.json();
  },

  async checkHealth() {
    const res = await fetch(`${this.baseUrl}/health`);
    if (!res.ok) throw new Error('服务异常');
    return res.json();
  },

  // 后续可扩展：获取设备列表、数据详情等
  // async getDevices() { ... }
  // async getDataByDevice(deviceId) { ... }
};
