<template>
  <div class="manage-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">数据管理</h1>
        <p class="page-desc">管理产品类型、功能点和安全产品的增删改查</p>
      </div>
    </header>

    <div class="page-content">
      <!-- Tabs -->
      <div class="tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          :class="['tab', { active: activeTab === tab.id }]"
          @click="activeTab = tab.id"
        >
          <span class="tab-icon" v-html="tab.icon"></span>
          <span class="tab-label">{{ tab.label }}</span>
          <span class="tab-count">{{ tabCounts[tab.id] || 0 }}</span>
        </button>
      </div>

      <!-- Product Types Tab -->
      <div v-show="activeTab === 'types'" class="tab-content">
        <div class="content-header">
          <div class="search-box">
            <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
            <input
              v-model="typeSearch"
              type="text"
              placeholder="搜索产品类型..."
              class="search-input"
            />
          </div>
          <button class="btn btn-primary" @click="openTypeModal()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            新增类型
          </button>
        </div>

        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>类型名称</th>
                <th>描述</th>
                <th>图标</th>
                <th>颜色</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in filteredTypes" :key="item.id">
                <td class="mono">{{ item.id }}</td>
                <td>
                  <div class="type-name-cell">
                    <span class="type-dot" :style="{ background: item.color }"></span>
                    {{ item.name }}
                  </div>
                </td>
                <td class="desc-cell">{{ item.description || '-' }}</td>
                <td class="icon-cell">{{ item.icon || '-' }}</td>
                <td>
                  <span class="color-swatch" :style="{ background: item.color }"></span>
                  {{ item.color }}
                </td>
                <td>
                  <div class="action-btns">
                    <button class="action-btn edit" @click="openTypeModal(item)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="deleteType(item)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="!filteredTypes.length" class="empty-state">
            <p>暂无数据</p>
          </div>
        </div>
      </div>

      <!-- Functions Tab -->
      <div v-show="activeTab === 'functions'" class="tab-content">
        <div class="content-header">
          <div class="search-box">
            <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
            <input
              v-model="funcSearch"
              type="text"
              placeholder="搜索功能点..."
              class="search-input"
            />
          </div>
          <div class="filter-group">
            <select v-model="funcCategoryFilter" class="filter-select">
              <option value="">全部分类</option>
              <option v-for="cat in funcCategories" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>
          <button class="btn btn-primary" @click="openFuncModal()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            新增功能
          </button>
        </div>

        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>功能名称</th>
                <th>分类</th>
                <th>描述</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in filteredFunctions" :key="item.id">
                <td class="mono">{{ item.id }}</td>
                <td class="name-cell">{{ item.name }}</td>
                <td>
                  <span class="category-tag">{{ item.category }}</span>
                </td>
                <td class="desc-cell">{{ item.description || '-' }}</td>
                <td>
                  <div class="action-btns">
                    <button class="action-btn edit" @click="openFuncModal(item)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="deleteFunc(item)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="!filteredFunctions.length" class="empty-state">
            <p>暂无数据</p>
          </div>
        </div>
      </div>

      <!-- Products Tab -->
      <div v-show="activeTab === 'products'" class="tab-content">
        <div class="content-header">
          <div class="search-box">
            <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
            <input
              v-model="prodSearch"
              type="text"
              placeholder="搜索产品..."
              class="search-input"
            />
          </div>
          <div class="filter-group">
            <select v-model="prodTypeFilter" class="filter-select">
              <option value="">全部类型</option>
              <option v-for="t in productTypes" :key="t.id" :value="t.id">{{ t.name }}</option>
            </select>
          </div>
          <button class="btn btn-primary" @click="openProdModal()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            新增产品
          </button>
        </div>

        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>产品名称</th>
                <th>品牌</th>
                <th>类型</th>
                <th>功能数量</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in filteredProducts" :key="item.id">
                <td class="mono">{{ item.id }}</td>
                <td class="name-cell">{{ item.name }}</td>
                <td>{{ item.brand || '-' }}</td>
                <td>
                  <span class="type-badge" :style="{ borderColor: getTypeColor(item.type_id), color: getTypeColor(item.type_id) }">
                    {{ getTypeName(item.type_id) }}
                  </span>
                </td>
                <td>
                  <span class="func-count">{{ item.function_ids?.length || 0 }}</span>
                </td>
                <td>
                  <div class="action-btns">
                    <button class="action-btn edit" @click="openProdModal(item)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="deleteProd(item)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="!filteredProducts.length" class="empty-state">
            <p>暂无数据</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Type Modal -->
    <div v-if="showTypeModal" class="modal-overlay" @click.self="closeTypeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingType ? '编辑类型' : '新增类型' }}</h3>
          <button class="modal-close" @click="closeTypeModal">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>类型名称 *</label>
            <input v-model="typeForm.name" type="text" class="form-input" placeholder="如：防火墙、入侵检测系统" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="typeForm.description" class="form-input form-textarea" placeholder="描述此产品类型的用途..."></textarea>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>图标标识</label>
              <input v-model="typeForm.icon" type="text" class="form-input" placeholder="如：firewall" />
            </div>
            <div class="form-group">
              <label>颜色</label>
              <div class="color-picker">
                <input v-model="typeForm.color" type="color" class="color-input" />
                <input v-model="typeForm.color" type="text" class="form-input color-text" placeholder="#ff6b6b" />
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeTypeModal">取消</button>
          <button class="btn btn-primary" @click="saveType">{{ editingType ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- Function Modal -->
    <div v-if="showFuncModal" class="modal-overlay" @click.self="closeFuncModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingFunc ? '编辑功能点' : '新增功能点' }}</h3>
          <button class="modal-close" @click="closeFuncModal">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>功能名称 *</label>
            <input v-model="funcForm.name" type="text" class="form-input" placeholder="如：入侵检测、访问控制" />
          </div>
          <div class="form-group">
            <label>分类 *</label>
            <select v-model="funcForm.category" class="form-input">
              <option value="">请选择分类</option>
              <option v-for="cat in funcCategories" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="funcForm.description" class="form-input form-textarea" placeholder="描述此功能点的具体能力..."></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeFuncModal">取消</button>
          <button class="btn btn-primary" @click="saveFunc">{{ editingFunc ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- Product Modal -->
    <div v-if="showProdModal" class="modal-overlay" @click.self="closeProdModal">
      <div class="modal modal-lg">
        <div class="modal-header">
          <h3>{{ editingProd ? '编辑产品' : '新增产品' }}</h3>
          <button class="modal-close" @click="closeProdModal">×</button>
        </div>
        <div class="modal-body">
          <div class="form-row">
            <div class="form-group">
              <label>产品名称 *</label>
              <input v-model="prodForm.name" type="text" class="form-input" placeholder="如：天擎防火墙" />
            </div>
            <div class="form-group">
              <label>品牌</label>
              <input v-model="prodForm.brand" type="text" class="form-input" placeholder="如：奇安信" />
            </div>
          </div>
          <div class="form-group">
            <label>产品类型 *</label>
            <select v-model="prodForm.type_id" class="form-input">
              <option value="">请选择类型</option>
              <option v-for="t in productTypes" :key="t.id" :value="t.id">{{ t.name }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>关联功能点</label>
            <div class="function-selector">
              <div class="function-category" v-for="cat in funcCategories" :key="cat">
                <div class="category-label">{{ cat }}</div>
                <div class="function-checks">
                  <label
                    v-for="func in getFunctionsByCategory(cat)"
                    :key="func.id"
                    class="function-check"
                  >
                    <input
                      type="checkbox"
                      :value="func.id"
                      v-model="prodForm.function_ids"
                    />
                    <span>{{ func.name }}</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeProdModal">取消</button>
          <button class="btn btn-primary" @click="saveProd">{{ editingProd ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="closeDeleteModal">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>确认删除</h3>
          <button class="modal-close" @click="closeDeleteModal">×</button>
        </div>
        <div class="modal-body">
          <p class="delete-message">{{ deleteConfirmText }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeDeleteModal">取消</button>
          <button class="btn btn-danger" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { productTypeApi, functionApi, productApi } from '../api'
import api from '../api'

const tabs = [
  {
    id: 'types',
    label: '产品类型',
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 9h6M9 15h6M9 12h6"/></svg>'
  },
  {
    id: 'functions',
    label: '功能点',
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 1v6M12 17v6M4.22 4.22l4.24 4.24M15.54 15.54l4.24 4.24M1 12h6M17 12h6M4.22 19.78l4.24-4.24M15.54 8.46l4.24-4.24"/></svg>'
  },
  {
    id: 'products',
    label: '安全产品',
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v2"/></svg>'
  }
]

// Data
const activeTab = ref('types')
const productTypes = ref([])
const functions = ref([])
const products = ref([])

// Search/Filter
const typeSearch = ref('')
const funcSearch = ref('')
const funcCategoryFilter = ref('')
const prodSearch = ref('')
const prodTypeFilter = ref('')

// Tab counts
const tabCounts = computed(() => ({
  types: productTypes.value.length,
  functions: functions.value.length,
  products: products.value.length
}))

// Categories
const funcCategories = computed(() => {
  const cats = new Set()
  functions.value.forEach(f => { if (f.category) cats.add(f.category) })
  return Array.from(cats).sort()
})

// Filtered data
const filteredTypes = computed(() => {
  if (!typeSearch.value) return productTypes.value
  const q = typeSearch.value.toLowerCase()
  return productTypes.value.filter(t =>
    t.name.toLowerCase().includes(q) ||
    (t.description && t.description.toLowerCase().includes(q))
  )
})

const filteredFunctions = computed(() => {
  let list = functions.value
  if (funcCategoryFilter.value) {
    list = list.filter(f => f.category === funcCategoryFilter.value)
  }
  if (!funcSearch.value) return list
  const q = funcSearch.value.toLowerCase()
  return list.filter(f =>
    f.name.toLowerCase().includes(q) ||
    (f.description && f.description.toLowerCase().includes(q))
  )
})

const filteredProducts = computed(() => {
  let list = products.value
  if (prodTypeFilter.value) {
    list = list.filter(p => p.type_id === prodTypeFilter.value)
  }
  if (!prodSearch.value) return list
  const q = prodSearch.value.toLowerCase()
  return list.filter(p =>
    p.name.toLowerCase().includes(q) ||
    (p.brand && p.brand.toLowerCase().includes(q))
  )
})

// Type Modal
const showTypeModal = ref(false)
const editingType = ref(null)
const typeForm = ref({ name: '', description: '', icon: '', color: '#00d4aa' })

// Func Modal
const showFuncModal = ref(false)
const editingFunc = ref(null)
const funcForm = ref({ name: '', category: '', description: '' })

// Prod Modal
const showProdModal = ref(false)
const editingProd = ref(null)
const prodForm = ref({ name: '', brand: '', type_id: '', function_ids: [] })

// Delete Modal
const showDeleteModal = ref(false)
const deleteTarget = ref(null)
const deleteCategory = ref('')
const deleteConfirmText = computed(() => {
  if (!deleteTarget.value) return ''
  const names = { type: '类型', function: '功能点', product: '产品' }
  return `确定要删除 "${deleteTarget.value.name}" 吗？此操作无法撤销。`
})

// Methods
function getTypeName(typeId) {
  const t = productTypes.value.find(pt => pt.id === typeId)
  return t ? t.name : '-'
}

function getTypeColor(typeId) {
  const t = productTypes.value.find(pt => pt.id === typeId)
  return t ? t.color : '#666'
}

function getFunctionsByCategory(category) {
  return functions.value.filter(f => f.category === category)
}

// Type operations
function openTypeModal(item = null) {
  editingType.value = item
  if (item) {
    typeForm.value = { ...item }
  } else {
    typeForm.value = { name: '', description: '', icon: '', color: '#00d4aa' }
  }
  showTypeModal.value = true
}

function closeTypeModal() {
  showTypeModal.value = false
  editingType.value = null
}

async function saveType() {
  if (!typeForm.value.name) {
    alert('请输入类型名称')
    return
  }
  try {
    if (editingType.value) {
      await productTypeApi.update(editingType.value.id, typeForm.value)
    } else {
      await productTypeApi.create(typeForm.value)
    }
    await loadData()
    closeTypeModal()
  } catch (err) {
    alert(err.response?.data?.message || '保存失败')
  }
}

// Function operations
function openFuncModal(item = null) {
  editingFunc.value = item
  if (item) {
    funcForm.value = { ...item }
  } else {
    funcForm.value = { name: '', category: '', description: '' }
  }
  showFuncModal.value = true
}

function closeFuncModal() {
  showFuncModal.value = false
  editingFunc.value = null
}

async function saveFunc() {
  if (!funcForm.value.name || !funcForm.value.category) {
    alert('请输入功能名称和分类')
    return
  }
  try {
    if (editingFunc.value) {
      await functionApi.update(editingFunc.value.id, funcForm.value)
    } else {
      await functionApi.create(funcForm.value)
    }
    await loadData()
    closeFuncModal()
  } catch (err) {
    alert(err.response?.data?.message || '保存失败')
  }
}

// Product operations
function openProdModal(item = null) {
  editingProd.value = item
  if (item) {
    prodForm.value = {
      name: item.name,
      brand: item.brand,
      type_id: item.type_id,
      function_ids: [...(item.function_ids || [])]
    }
  } else {
    prodForm.value = { name: '', brand: '', type_id: '', function_ids: [] }
  }
  showProdModal.value = true
}

function closeProdModal() {
  showProdModal.value = false
  editingProd.value = null
}

async function saveProd() {
  if (!prodForm.value.name || !prodForm.value.type_id) {
    alert('请输入产品名称和选择类型')
    return
  }
  try {
    const payload = {
      name: prodForm.value.name,
      brand: prodForm.value.brand,
      type_id: Number(prodForm.value.type_id),
      function_ids: prodForm.value.function_ids.map(Number)
    }
    if (editingProd.value) {
      await productApi.update(editingProd.value.id, payload)
    } else {
      await productApi.create(payload)
    }
    await loadData()
    closeProdModal()
  } catch (err) {
    alert(err.response?.data?.message || '保存失败')
  }
}

// Delete operations
function deleteType(item) {
  deleteTarget.value = item
  deleteCategory.value = 'type'
  showDeleteModal.value = true
}

function deleteFunc(item) {
  deleteTarget.value = item
  deleteCategory.value = 'function'
  showDeleteModal.value = true
}

function deleteProd(item) {
  deleteTarget.value = item
  deleteCategory.value = 'product'
  showDeleteModal.value = true
}

function closeDeleteModal() {
  showDeleteModal.value = false
  deleteTarget.value = null
  deleteCategory.value = ''
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    const deleteFunctions = {
      type: productTypeApi.delete,
      function: functionApi.delete,
      product: productApi.delete
    }
    const deleteFn = deleteFunctions[deleteCategory.value]
    await deleteFn(deleteTarget.value.id)
    await loadData()
    closeDeleteModal()
  } catch (err) {
    alert(err.response?.data?.message || '删除失败')
  }
}

// Load data
async function loadData() {
  try {
    const [typesRes, funcsRes, prodsRes] = await Promise.all([
      productApi.list(),
      functionApi.list(),
      productApi.listProducts()
    ])
    productTypes.value = typesRes || []
    functions.value = funcsRes || []
    products.value = prodsRes || []
  } catch (err) {
    console.error('Failed to load data:', err)
  }
}

onMounted(loadData)
</script>

<style scoped>
/* Layout */
.manage-page {
  min-height: 100vh;
}

.page-header {
  padding: 32px 40px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
}

.page-content {
  padding: 32px 40px;
}

/* Tabs */
.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  padding: 6px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-lg);
  width: fit-content;
}

.tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
}

.tab.active {
  background: var(--bg-card);
  color: var(--accent-primary);
  box-shadow: var(--shadow-sm);
}

.tab-icon {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab-icon :deep(svg) {
  width: 18px;
  height: 18px;
}

.tab-count {
  padding: 2px 8px;
  background: var(--bg-primary);
  border-radius: 10px;
  font-size: 12px;
  font-family: var(--font-mono);
}

.tab.active .tab-count {
  background: rgba(0, 212, 170, 0.15);
  color: var(--accent-primary);
}

/* Content Header */
.content-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  color: var(--text-tertiary);
}

.search-input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 14px;
}

.search-input:focus {
  outline: none;
  border-color: var(--accent-primary);
}

.filter-select {
  min-width: 150px;
}

/* Table */
.table-wrapper {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.data-table th {
  background: var(--bg-tertiary);
  font-weight: 600;
  font-size: 13px;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table tr:hover td {
  background: var(--bg-hover);
}

.mono {
  font-family: var(--font-mono);
  color: var(--text-tertiary);
}

.name-cell {
  font-weight: 500;
  color: var(--text-primary);
}

.desc-cell {
  color: var(--text-secondary);
  font-size: 13px;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.icon-cell {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-tertiary);
}

.type-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.type-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.color-swatch {
  display: inline-block;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  vertical-align: middle;
  margin-right: 6px;
}

.category-tag {
  display: inline-block;
  padding: 4px 10px;
  background: rgba(88, 166, 255, 0.1);
  color: var(--info);
  border-radius: 4px;
  font-size: 12px;
}

.type-badge {
  display: inline-block;
  padding: 4px 10px;
  border: 1px solid;
  border-radius: 4px;
  font-size: 12px;
  background: transparent;
}

.func-count {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--accent-primary);
}

/* Action Buttons */
.action-btns {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn:hover {
  border-color: var(--accent-primary);
  color: var(--accent-primary);
}

.action-btn.delete:hover {
  border-color: var(--danger);
  color: var(--danger);
  background: rgba(248, 81, 73, 0.1);
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal {
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  animation: slideUp 0.3s ease;
}

.modal-lg {
  max-width: 700px;
}

.modal-sm {
  max-width: 400px;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--text-tertiary);
  font-size: 24px;
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
}

.modal-close:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.modal-body {
  padding: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
}

/* Form */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 14px;
}

.form-input:focus {
  outline: none;
  border-color: var(--accent-primary);
}

.form-textarea {
  min-height: 80px;
  resize: vertical;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.color-picker {
  display: flex;
  gap: 8px;
  align-items: center;
}

.color-input {
  width: 40px;
  height: 40px;
  padding: 2px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.color-text {
  flex: 1;
}

/* Function Selector */
.function-selector {
  max-height: 300px;
  overflow-y: auto;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.function-category {
  margin-bottom: 16px;
}

.function-category:last-child {
  margin-bottom: 0;
}

.category-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--border-color);
}

.function-checks {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.function-check {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: all 0.2s ease;
}

.function-check:hover {
  border-color: var(--accent-primary);
}

.function-check:has(input:checked) {
  background: rgba(0, 212, 170, 0.1);
  border-color: var(--accent-primary);
  color: var(--accent-primary);
}

.function-check input {
  display: none;
}

/* Delete */
.delete-message {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.btn-danger {
  background: var(--danger);
  color: white;
}

.btn-danger:hover {
  background: #dc3545;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn svg {
  width: 18px;
  height: 18px;
}
</style>
