# 信创网络安全能力协同分析系统 - 前端

## 技术栈

- Vue 3 + Composition API
- Vue Router 4
- Pinia (状态管理)
- Axios (HTTP 客户端)
- Vis-Network (网络拓扑图可视化)
- Vite (构建工具)

## 快速开始

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器将运行在 `http://localhost:3000`，并代理 API 请求到 `http://localhost:8080`。

## 项目结构

```
frontend/
├── public/
│   └── favicon.svg
├── src/
│   ├── api/          # API 调用封装
│   ├── components/   # 公共组件
│   ├── router/       # 路由配置
│   ├── styles/       # 全局样式
│   ├── views/        # 页面组件
│   ├── App.vue
│   └── main.js
├── index.html
├── package.json
└── vite.config.js
```

## 页面说明

### 概览 (/).
展示系统整体状态：拓扑数量、产品数量、功能覆盖统计、拓扑列表预览。

### 网络拓扑 (/topology)
- 选择拓扑查看网络结构
- 支持节点点击查看详情（产品信息、节点属性、安全功能）
- 攻击路径分析面板
- 节点颜色区分类型（硬件/软件/OS/服务）

### 能力分析 (/analysis)
- 选择拓扑进行安全能力分析
- 覆盖率/冗余率可视化
- 缺失能力列表
- 风险评估（覆盖率缺口、冗余率、单点承载）

### 优化建议 (/suggest/:id)
- 双策略切换（最小变动/最小规模）
- 新增/移除产品建议
- 优化效果量化展示

## 构建

```bash
npm run build   # 构建生产版本
npm run preview # 预览生产构建
```
