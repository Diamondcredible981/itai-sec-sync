<template>
  <div class="suggest-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">优化建议</h1>
        <p class="page-desc">智能分析并生成最优产品配置方案</p>
      </div>
      <div class="strategy-toggle">
        <button
          class="strategy-btn"
          :class="{ active: strategy === 'min-change' }"
          @click="strategy = 'min-change'"
        >
          最小变动
        </button>
        <button
          class="strategy-btn"
          :class="{ active: strategy === 'min-size' }"
          @click="strategy = 'min-size'"
        >
          最小规模
        </button>
      </div>
    </header>

    <div class="page-content" v-if="displayState === 'result'">
      <div class="result-header">
        <div class="result-status" :class="{ success: result.success }">
          <svg v-if="result.success" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
          <span>{{ result.message || (result.success ? '方案生成成功' : '方案生成失败') }}</span>
        </div>
        <div class="operation-count">
          <span class="count">{{ result.total_operations }}</span>
          <span class="label">项调整</span>
        </div>
      </div>

      <div class="explanation-card" v-if="result.explanation">
        <h3>优化效果</h3>
        <div class="metrics-grid">
          <div class="metric">
            <div class="metric-change" :class="getDeltaClass(result.explanation.improvement.coverage_rate_delta)">
              {{ formatDelta(result.explanation.improvement.coverage_rate_delta) }}
            </div>
            <div class="metric-label">覆盖率</div>
            <div class="metric-values">
              <span>{{ result.explanation.current.coverage_rate }}%</span>
              →
              <span>{{ result.explanation.target.coverage_rate }}%</span>
            </div>
          </div>
          <div class="metric">
            <div class="metric-change" :class="getDeltaClass(-result.explanation.improvement.redundancy_rate_delta)">
              {{ formatDelta(result.explanation.improvement.redundancy_rate_delta) }}
            </div>
            <div class="metric-label">冗余率</div>
            <div class="metric-values">
              <span>{{ result.explanation.current.redundancy_rate }}%</span>
              →
              <span>{{ result.explanation.target.redundancy_rate }}%</span>
            </div>
          </div>
          <div class="metric">
            <div class="metric-change" :class="getDeltaClass(result.explanation.improvement.risk_score_reduction)">
              {{ formatDelta(-result.explanation.improvement.risk_score_reduction) }}
            </div>
            <div class="metric-label">风险评分</div>
            <div class="metric-values">
              <span>{{ result.explanation.current.risk_score }}</span>
              →
              <span>{{ result.explanation.target.risk_score }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="operations-section" v-if="result.success">
        <div class="operations-grid">
          <div class="card add-card" v-if="result.add_operations?.length">
            <div class="card-header">
              <h3 class="card-title">新增产品</h3>
              <span class="badge badge-success">{{ result.add_operations.length }}</span>
            </div>
            <div class="operation-list">
              <div v-for="op in result.add_operations" :key="op.product_id" class="operation-item">
                <div class="op-icon add">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                </div>
                <div class="op-info">
                  <div class="op-name">{{ op.product.name }}</div>
                  <div class="op-brand">{{ op.product.brand }}</div>
                  <div class="op-reason">{{ op.reason }}</div>
                </div>
                <div class="op-type-badge">{{ getTypeName(op.product.type_id) }}</div>
              </div>
            </div>
          </div>

          <div class="card remove-card" v-if="result.remove_operations?.length">
            <div class="card-header">
              <h3 class="card-title">移除产品</h3>
              <span class="badge badge-danger">{{ result.remove_operations.length }}</span>
            </div>
            <div class="operation-list">
              <div v-for="op in result.remove_operations" :key="op.product_id" class="operation-item">
                <div class="op-icon remove">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                </div>
                <div class="op-info">
                  <div class="op-name">{{ op.product.name }}</div>
                  <div class="op-brand">{{ op.product.brand }}</div>
                  <div class="op-reason">{{ op.reason }}</div>
                </div>
                <div class="op-type-badge">{{ getTypeName(op.product.type_id) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="back-section">
        <router-link to="/analysis" class="btn btn-secondary">
          ← 返回分析
        </router-link>
      </div>
    </div>

    <div class="empty-state" v-if="displayState === 'no-topo'">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M12 2L2 7l10 5 10-5-10-5z"/>
        <path d="M2 17l10 5 10-5"/>
        <path d="M2 12l10 5 10-5"/>
      </svg>
      <p>请先从「能力分析」页面选择一个拓扑</p>
      <router-link to="/analysis" class="btn btn-primary" style="margin-top: 16px;">
        前往能力分析
      </router-link>
    </div>

    <div class="loading" v-if="displayState === 'loading'">
      <div class="spinner"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { suggestApi, productApi } from '../api'

const route = useRoute()
const loading = ref(false)
const strategy = ref('min-change')
const result = ref(null)
const productTypes = ref([])
const hasTopoId = ref(false)

// 显示状态计算属性
const displayState = computed(() => {
  if (loading.value) return 'loading'
  if (!hasTopoId.value) return 'no-topo'
  if (result.value) return 'result'
  return 'loading'
})

function getTypeName(typeId) {
  const t = productTypes.value.find(t => t.id === typeId)
  return t ? t.name : '未知'
}

function getDeltaClass(delta) {
  if (delta > 0) return 'positive'
  if (delta < 0) return 'negative'
  return 'neutral'
}

function formatDelta(delta) {
  if (delta > 0) return `+${delta.toFixed(1)}`
  if (delta < 0) return delta.toFixed(1)
  return '0'
}

async function loadProductTypes() {
  try {
    const res = await productApi.list()
    productTypes.value = res || []
  } catch (err) {
    console.error('Failed to load product types:', err)
  }
}

async function loadSuggestion() {
  const id = route.params.id
  hasTopoId.value = !!id

  if (!id) {
    return
  }

  loading.value = true
  try {
    const res = await suggestApi.get(id, strategy.value)
    result.value = res
  } catch (err) {
    console.error('Failed to load suggestion:', err)
  } finally {
    loading.value = false
  }
}

watch(strategy, () => {
  if (route.params.id) {
    loadSuggestion()
  }
})

watch(() => route.params.id, (newId) => {
  hasTopoId.value = !!newId
  if (newId) {
    loadSuggestion()
  } else {
    result.value = null
  }
})

onMounted(async () => {
  await loadProductTypes()
  await loadSuggestion()
})
</script>

<style scoped>
.strategy-toggle {
  display: flex;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  padding: 4px;
}

.strategy-btn {
  padding: 10px 20px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
}

.strategy-btn.active {
  background: var(--accent-gradient);
  color: white;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  margin-bottom: 24px;
}

.result-status {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 500;
  color: var(--danger);
}

.result-status.success {
  color: var(--success);
}

.result-status svg {
  width: 24px;
  height: 24px;
}

.operation-count {
  display: flex;
  align-items: baseline;
  gap: 6px;
  text-align: right;
}

.operation-count .count {
  font-size: 28px;
  font-weight: 700;
  font-family: var(--font-mono);
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

.operation-count .label {
  font-size: 13px;
  color: var(--text-tertiary);
}

.explanation-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 32px;
  margin-bottom: 32px;
}

.explanation-card h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 24px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
}

@media (max-width: 768px) {
  .metrics-grid {
    grid-template-columns: 1fr;
  }
}

.metric {
  text-align: center;
}

.metric-change {
  font-size: 28px;
  font-weight: 700;
  font-family: var(--font-mono);
  margin-bottom: 8px;
}

.metric-change.positive {
  color: var(--success);
}

.metric-change.negative {
  color: var(--danger);
}

.metric-change.neutral {
  color: var(--text-tertiary);
}

.metric-label {
  font-size: 13px;
  color: var(--text-tertiary);
  margin-bottom: 8px;
}

.metric-values {
  font-size: 14px;
  color: var(--text-secondary);
  font-family: var(--font-mono);
}

.operations-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

@media (max-width: 900px) {
  .operations-grid {
    grid-template-columns: 1fr;
  }
}

.operation-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.operation-item {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.op-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}

.op-icon.add {
  background: rgba(63, 185, 80, 0.1);
  color: var(--success);
}

.op-icon.remove {
  background: rgba(248, 81, 73, 0.1);
  color: var(--danger);
}

.op-icon svg {
  width: 18px;
  height: 18px;
}

.op-info {
  flex: 1;
}

.op-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.op-brand {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
}

.op-reason {
  font-size: 13px;
  color: var(--text-secondary);
}

.op-type-badge {
  font-size: 11px;
  color: var(--accent-primary);
  padding: 4px 10px;
  background: rgba(0, 212, 170, 0.1);
  border-radius: 12px;
}

.back-section {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

.operations-section {
  width: 100%;
}

.operations-grid {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

@media (max-width: 900px) {
  .operations-grid {
    grid-template-columns: 1fr;
  }
}

.operations-grid .card {
  width: 100%;
  min-width: 0;
}

.operations-grid h3.card-title {
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
  white-space: nowrap;
}

.operations-grid h3.card-title::before {
  content: '';
  display: inline-block;
  width: 4px;
  height: 16px;
  background: var(--accent-gradient);
  border-radius: 2px;
  margin-right: 10px;
  flex-shrink: 0;
}

.operations-grid .card-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.operations-grid .badge {
  flex-shrink: 0;
}
</style>
