<template>
  <div class="topology-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">网络拓扑</h1>
        <p class="page-desc">可视化展示网络节点结构与连接关系</p>
      </div>
      <div class="header-actions">
        <select v-model="selectedTopoId" @change="loadTopology" class="topo-select">
          <option value="">选择拓扑</option>
          <option v-for="t in topologies" :key="t.id" :value="t.id">
            {{ t.name }}
          </option>
        </select>
        <button class="btn btn-secondary" @click="refreshTopology">刷新</button>
      </div>
    </header>

    <div class="topology-content" v-if="currentTopo">
      <div class="topology-toolbar">
        <div class="toolbar-left">
          <span class="topo-name">{{ currentTopo.name }}</span>
          <span class="node-count">{{ nodes.length }} 节点</span>
          <span class="edge-count">{{ edges.length }} 边</span>
        </div>
        <div class="toolbar-right">
          <button
            class="btn btn-secondary"
            :class="{ active: showAttackPath }"
            @click="showAttackPath = !showAttackPath"
          >
            攻击路径
          </button>
        </div>
      </div>

      <div class="topology-container">
        <div class="graph-wrapper">
          <div ref="graphContainer" class="graph-container"></div>
        </div>

        <transition name="slide">
          <div class="node-panel" v-if="selectedNode">
            <div class="panel-header">
              <h3>节点详情</h3>
              <button class="close-btn" @click="selectedNode = null">×</button>
            </div>
            <div class="panel-content">
              <div class="node-header">
                <div class="node-icon" :style="{ background: getNodeColor(selectedNode) + '20', color: getNodeColor(selectedNode) }">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="3" width="18" height="18" rx="2"/>
                    <path d="M9 9h6M9 15h6"/>
                  </svg>
                </div>
                <div class="node-title">
                  <div class="node-name">{{ selectedNode.label }}</div>
                  <div class="node-type">{{ selectedNode.data?.node_type || selectedNode.type }}</div>
                </div>
              </div>

              <div class="info-section" v-if="selectedNode.product">
                <h4>产品信息</h4>
                <div class="info-grid">
                  <div class="info-item">
                    <span class="info-label">产品名称</span>
                    <span class="info-value">{{ selectedNode.product.name }}</span>
                  </div>
                  <div class="info-item">
                    <span class="info-label">品牌</span>
                    <span class="info-value">{{ selectedNode.product.brand }}</span>
                  </div>
                  <div class="info-item">
                    <span class="info-label">类型</span>
                    <span class="info-value">{{ selectedNode.product.type }}</span>
                  </div>
                </div>
              </div>

              <div class="info-section" v-if="selectedNode.data">
                <h4>节点属性</h4>
                <div class="info-grid">
                  <div class="info-item">
                    <span class="info-label">Zone</span>
                    <span class="info-value">{{ selectedNode.data.zone }}</span>
                  </div>
                  <div class="info-item">
                    <span class="info-label">Layer</span>
                    <span class="info-value">{{ selectedNode.data.layer }}</span>
                  </div>
                  <div class="info-item">
                    <span class="info-label">关键性</span>
                    <span class="info-value">{{ selectedNode.data.criticality }}</span>
                  </div>
                </div>
              </div>

              <div class="info-section" v-if="selectedNode.product?.function_ids?.length">
                <h4>安全功能</h4>
                <div class="function-tags">
                  <span
                    v-for="fid in selectedNode.product.function_ids"
                    :key="fid"
                    class="function-tag"
                  >
                    {{ getFunctionName(fid) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </transition>
      </div>

      <div class="attack-path-panel" v-if="showAttackPath && attackPath">
        <div class="path-header">
          <h3>攻击路径分析</h3>
          <button class="close-btn" @click="showAttackPath = false">×</button>
        </div>
        <div class="path-summary">
          <div class="path-stat">
            <span class="stat-label">路径长度</span>
            <span class="stat-value">{{ attackPath.path_length }}</span>
          </div>
          <div class="path-stat">
            <span class="stat-label">整体风险</span>
            <span class="stat-value" :class="getRiskClass(attackPath.risk_level)">
              {{ attackPath.risk_level }}
            </span>
          </div>
          <div class="path-stat">
            <span class="stat-label">风险评分</span>
            <span class="stat-value">{{ attackPath.overall_risk }}</span>
          </div>
        </div>
        <div class="path-mitigations" v-if="attackPath.mitigations?.length">
          <h4>缓解建议</h4>
          <ul>
            <li v-for="(m, i) in attackPath.mitigations" :key="i">{{ m }}</li>
          </ul>
        </div>
      </div>
    </div>

    <div class="empty-state" v-else-if="!loading">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="6" cy="6" r="3"/>
        <circle cx="18" cy="6" r="3"/>
        <circle cx="6" cy="18" r="3"/>
        <circle cx="18" cy="18" r="3"/>
        <path d="M6 9v6M18 9v6M9 6h6M9 18h6"/>
      </svg>
      <p>请选择一个拓扑以开始分析</p>
    </div>

    <div class="loading" v-if="loading">
      <div class="spinner"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { topoApi, functionApi } from '../api'
import { Network, DataSet } from 'vis-network/standalone'
import 'vis-network/styles/vis-network.css'

const route = useRoute()
const graphContainer = ref(null)
let network = null

const loading = ref(false)
const topologies = ref([])
const selectedTopoId = ref('')
const currentTopo = ref(null)
const nodes = ref([])
const edges = ref([])
const selectedNode = ref(null)
const showAttackPath = ref(false)
const attackPath = ref(null)
const functions = ref([])

const nodeColors = {
  hardware: '#00d4aa',
  software: '#58a6ff',
  os: '#a371f7',
  service: '#f0883e'
}

function getNodeColor(node) {
  if (node.category) return node.category === 'firewall' ? '#ff6b6b' : nodeColors[node.type] || '#8b949e'
  return nodeColors[node.type] || '#8b949e'
}

function getFunctionName(id) {
  const f = functions.value.find(f => f.id === id)
  return f ? f.name : `功能${id}`
}

function getRiskClass(level) {
  const map = { critical: 'risk-critical', high: 'risk-high', medium: 'risk-medium', low: 'risk-low' }
  return map[level] || ''
}

async function loadTopoList() {
  try {
    const res = await topoApi.list({ mode: 'summary' })
    topologies.value = res.items || res || []
  } catch (err) {
    console.error('Failed to load topologies:', err)
  }
}

async function loadTopology() {
  if (!selectedTopoId.value) {
    currentTopo.value = null
    nodes.value = []
    edges.value = []
    return
  }

  loading.value = true
  try {
    const visRes = await topoApi.visualization(selectedTopoId.value)
    const detailRes = await topoApi.get(selectedTopoId.value)

    currentTopo.value = detailRes
    nodes.value = visRes.nodes || []
    edges.value = visRes.edges || []

    await nextTick()
    initGraph()
    loadAttackPath()
  } catch (err) {
    console.error('Failed to load topology:', err)
  } finally {
    loading.value = false
  }
}

async function loadFunctions() {
  try {
    const res = await functionApi.list()
    functions.value = res || []
  } catch (err) {
    console.error('Failed to load functions:', err)
  }
}

async function loadAttackPath() {
  if (!selectedTopoId.value) return
  try {
    const res = await topoApi.attackPath(selectedTopoId.value, {})
    attackPath.value = res
  } catch (err) {
    console.error('Failed to load attack path:', err)
  }
}

function initGraph() {
  if (!graphContainer.value || !nodes.value.length) return

  if (network) {
    network.destroy()
    network = null
  }

  const nodesData = new DataSet(nodes.value.map(n => ({
    id: n.id,
    label: n.label,
    title: `${n.label}\n类型: ${n.type || n.data?.node_type || '未知'}`,
    color: {
      background: getNodeColor(n) + '30',
      border: getNodeColor(n),
      highlight: { background: getNodeColor(n) + '50', border: getNodeColor(n) }
    },
    font: { color: '#e6edf3', size: 14 },
    shape: n.product ? 'box' : 'ellipse',
    size: n.product ? 30 : 20,
    data: n.data || {},
    product: n.product
  })))

  const edgesData = new DataSet(edges.value.map(e => ({
    id: e.id,
    from: e.source,
    to: e.target,
    color: { color: '#4a5568', highlight: '#00d4aa' },
    width: e.weight || 1,
    arrows: e.direction === 'bi' ? 'to;from' : 'to',
    title: `风险: ${e.risk || 0}`
  })))

  const options = {
    nodes: {
      borderWidth: 2,
      shadow: true
    },
    edges: {
      smooth: { type: 'continuous' }
    },
    physics: {
      enabled: true,
      solver: 'repulsion',
      repulsion: { nodeDistance: 200 }
    },
    interaction: {
      hover: true,
      tooltipDelay: 200
    }
  }

  network = new Network(graphContainer.value, { nodes: nodesData, edges: edgesData }, options)

  network.on('click', (params) => {
    if (params.nodes.length > 0) {
      const nodeId = params.nodes[0]
      const node = nodes.value.find(n => n.id === nodeId)
      if (node) {
        selectedNode.value = node
      }
    }
  })
}

function refreshTopology() {
  loadTopology()
}

watch(() => route.params.id, (newId) => {
  if (newId) {
    selectedTopoId.value = newId
    loadTopology()
  }
}, { immediate: true })

onMounted(async () => {
  await loadTopoList()
  await loadFunctions()
  if (route.params.id) {
    selectedTopoId.value = route.params.id
    await loadTopology()
  }
})

onUnmounted(() => {
  if (network) {
    network.destroy()
    network = null
  }
})
</script>

<style scoped>
.topology-page {
  min-height: 100vh;
}

.header-left {
  flex: 1;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.topo-select {
  min-width: 200px;
}

.topology-content {
  padding: 24px 40px;
}

.topology-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  margin-bottom: 20px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.topo-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.node-count,
.edge-count {
  font-size: 13px;
  color: var(--text-tertiary);
  font-family: var(--font-mono);
}

.node-count::before,
.edge-count::before {
  content: '';
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 6px;
}

.node-count::before {
  background: var(--accent-primary);
}

.edge-count::before {
  background: var(--info);
}

.btn.active {
  background: var(--accent-gradient);
  color: white;
  border-color: transparent;
}

.topology-container {
  display: flex;
  gap: 20px;
  min-height: 600px;
}

.graph-wrapper {
  flex: 1;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.graph-container {
  width: 100%;
  height: 600px;
}

.node-panel {
  width: 360px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.panel-header h3 {
  font-size: 16px;
  font-weight: 600;
}

.close-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--text-tertiary);
  font-size: 20px;
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.node-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
}

.node-icon svg {
  width: 28px;
  height: 28px;
}

.node-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.node-type {
  font-size: 13px;
  color: var(--text-tertiary);
  text-transform: capitalize;
}

.info-section {
  margin-bottom: 24px;
}

.info-section h4 {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
}

.info-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 14px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
}

.info-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.info-value {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.function-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.function-tag {
  padding: 6px 12px;
  background: rgba(0, 212, 170, 0.1);
  color: var(--accent-primary);
  font-size: 12px;
  border-radius: 20px;
  font-weight: 500;
}

.attack-path-panel {
  margin-top: 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
}

.path-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.path-header h3 {
  font-size: 16px;
  font-weight: 600;
}

.path-summary {
  display: flex;
  gap: 32px;
  margin-bottom: 20px;
}

.path-stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-tertiary);
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  font-family: var(--font-mono);
}

.risk-critical { color: var(--danger); }
.risk-high { color: #f0883e; }
.risk-medium { color: var(--warning); }
.risk-low { color: var(--success); }

.path-mitigations h4 {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.path-mitigations ul {
  list-style: none;
  padding: 0;
}

.path-mitigations li {
  position: relative;
  padding-left: 20px;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-primary);
}

.path-mitigations li::before {
  content: '→';
  position: absolute;
  left: 0;
  color: var(--accent-primary);
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(20px);
  opacity: 0;
}
</style>
