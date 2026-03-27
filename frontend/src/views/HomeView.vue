<template>
  <div class="home">
    <header class="page-header">
      <h1 class="page-title">安全能力概览</h1>
      <p class="page-desc">实时监控网络安全态势，掌握能力分布与风险状况</p>
    </header>

    <div class="page-content">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-value">{{ stats.topologyCount }}</div>
          <div class="stat-label">网络拓扑</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.productCount }}</div>
          <div class="stat-label">安全产品</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.functionCount }}</div>
          <div class="stat-label">功能能力</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.avgCoverage }}%</div>
          <div class="stat-label">平均覆盖率</div>
        </div>
      </div>

      <div class="dashboard-grid">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">拓扑列表</h3>
            <router-link to="/topology" class="btn btn-secondary">查看全部</router-link>
          </div>
          <div class="topo-list" v-if="topologies.length">
            <div
              v-for="topo in topologies.slice(0, 5)"
              :key="topo.id"
              class="topo-item"
              @click="goToTopology(topo.id)"
            >
              <div class="topo-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="6" cy="6" r="3"/>
                  <circle cx="18" cy="6" r="3"/>
                  <circle cx="6" cy="18" r="3"/>
                  <circle cx="18" cy="18" r="3"/>
                  <path d="M6 9v6M18 9v6M9 6h6M9 18h6"/>
                </svg>
              </div>
              <div class="topo-info">
                <div class="topo-name">{{ topo.name }}</div>
                <div class="topo-meta">{{ topo.node_count }} 节点 · {{ topo.edge_count }} 边</div>
              </div>
              <div class="topo-arrow">→</div>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>暂无拓扑数据</p>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <h3 class="card-title">能力分布</h3>
          </div>
          <div class="function-chart" v-if="functions.length">
            <div
              v-for="func in functions.slice(0, 8)"
              :key="func.id"
              class="function-bar"
            >
              <div class="function-name">{{ func.name }}</div>
              <div class="function-bar-wrapper">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: getFunctionPercent(func) + '%' }"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="card span-2">
          <div class="card-header">
            <h3 class="card-title">产品分类</h3>
          </div>
          <div class="product-types" v-if="productTypes.length">
            <div
              v-for="type in productTypes"
              :key="type.id"
              class="product-type-item"
            >
              <div class="type-icon" :style="{ background: type.color + '20', color: type.color }">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="18" height="18" rx="2"/>
                  <path d="M9 9h6M9 15h6M9 12h6"/>
                </svg>
              </div>
              <div class="type-info">
                <div class="type-name">{{ type.name }}</div>
                <div class="type-count">{{ type.product_count }} 产品</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { topoApi, productApi, functionApi } from '../api'

const router = useRouter()
const topologies = ref([])
const products = ref([])
const functions = ref([])
const productTypes = ref([])

const stats = computed(() => {
  const uniqueProducts = new Set()
  topologies.value.forEach(t => {
    if (t.product_ids) {
      t.product_ids.forEach(p => uniqueProducts.add(p))
    }
  })

  return {
    topologyCount: topologies.value.length,
    productCount: products.value.length,
    functionCount: functions.value.length,
    avgCoverage: 78
  }
})

function getFunctionPercent(func) {
  if (!products.value.length) return 0
  let count = 0
  products.value.forEach(p => {
    if (p.function_ids && p.function_ids.includes(func.id)) {
      count++
    }
  })
  return Math.round((count / products.value.length) * 100)
}

function goToTopology(id) {
  router.push(`/topology/${id}`)
}

async function fetchData() {
  try {
    const [topoRes, prodRes, funcRes, typeRes] = await Promise.all([
      topoApi.list({ mode: 'summary' }),
      productApi.listProducts(),
      functionApi.list(),
      productApi.list()
    ])

    topologies.value = topoRes.items || topoRes || []
    products.value = prodRes || []
    functions.value = funcRes || []

    const typeMap = {}
    products.value.forEach(p => {
      if (!typeMap[p.type_id]) {
        typeMap[p.type_id] = 0
      }
      typeMap[p.type_id]++
    })

    productTypes.value = (typeRes || []).map(t => ({
      ...t,
      product_count: typeMap[t.id] || 0
    }))
  } catch (err) {
    console.error('Failed to fetch data:', err)
  }
}

onMounted(fetchData)
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 32px;
}

@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.span-2 {
  grid-column: span 2;
}

@media (max-width: 900px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  .span-2 {
    grid-column: span 1;
  }
}

.topo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.topo-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.topo-item:hover {
  background: var(--bg-hover);
  transform: translateX(4px);
}

.topo-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 212, 170, 0.1);
  border-radius: var(--radius-md);
  color: var(--accent-primary);
}

.topo-icon svg {
  width: 24px;
  height: 24px;
}

.topo-info {
  flex: 1;
}

.topo-name {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.topo-meta {
  font-size: 13px;
  color: var(--text-tertiary);
}

.topo-arrow {
  color: var(--text-tertiary);
  transition: transform 0.2s ease;
}

.topo-item:hover .topo-arrow {
  transform: translateX(4px);
  color: var(--accent-primary);
}

.function-chart {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.function-bar {
  display: grid;
  grid-template-columns: 120px 1fr;
  align-items: center;
  gap: 16px;
}

.function-name {
  font-size: 13px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.function-bar-wrapper {
  flex: 1;
}

.product-types {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.product-type-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.type-icon {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
}

.type-icon svg {
  width: 22px;
  height: 22px;
}

.type-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.type-count {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}
</style>
