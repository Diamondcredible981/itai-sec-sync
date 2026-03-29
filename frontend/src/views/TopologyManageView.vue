<template>
  <div class="manage-page">
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">拓扑管理</h1>
        <p class="page-desc">管理网络拓扑的增删改查，包括节点和边的配置</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary" @click="openTopoModal()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          新增拓扑
        </button>
      </div>
    </header>

    <div class="page-content">
      <!-- Search -->
      <div class="content-header">
        <div class="search-box">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索拓扑名称..."
            class="search-input"
          />
        </div>
      </div>

      <!-- Topology List -->
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>拓扑名称</th>
              <th>节点数</th>
              <th>边数</th>
              <th>设备数</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredTopologies" :key="item.id">
              <td class="mono">{{ item.id }}</td>
              <td class="name-cell">{{ item.name }}</td>
              <td>
                <span class="stat-badge nodes">{{ item.node_count }}</span>
              </td>
              <td>
                <span class="stat-badge edges">{{ item.edge_count }}</span>
              </td>
              <td>
                <span class="stat-badge products">{{ item.product_count }}</span>
              </td>
              <td>
                <div class="action-btns">
                  <button class="action-btn view" @click="openViewModal(item)" title="查看">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                      <circle cx="12" cy="12" r="3"/>
                    </svg>
                  </button>
                  <button class="action-btn edit" @click="openTopoModal(item)" title="编辑">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button class="action-btn graphic" @click="openGraphicModal(item)" title="图形编辑">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="3" y="3" width="18" height="18" rx="2"/>
                      <circle cx="8" cy="8" r="2"/>
                      <circle cx="16" cy="8" r="2"/>
                      <circle cx="8" cy="16" r="2"/>
                      <circle cx="16" cy="16" r="2"/>
                      <path d="M10 8h4M10 16h4M8 10v4M16 10v4"/>
                    </svg>
                  </button>
                  <button class="action-btn copy" @click="copyTopo(item)" title="复制">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                    </svg>
                  </button>
                  <button class="action-btn delete" @click="deleteTopo(item)" title="删除">
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
        <div v-if="!filteredTopologies.length" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="6" cy="6" r="3"/>
            <circle cx="18" cy="6" r="3"/>
            <circle cx="6" cy="18" r="3"/>
            <circle cx="18" cy="18" r="3"/>
            <path d="M6 9v6M18 9v6M9 6h6M9 18h6"/>
          </svg>
          <p>暂无拓扑数据</p>
        </div>
      </div>
    </div>

    <!-- Topology Form Modal (表格模式) -->
    <div v-if="showTopoModal" class="modal-overlay" @click.self="closeTopoModal">
      <div class="modal modal-xl">
        <div class="modal-header">
          <h3>{{ editingTopo ? '编辑拓扑' : '新增拓扑' }}</h3>
          <button class="modal-close" @click="closeTopoModal">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>拓扑名称 *</label>
            <input v-model="topoForm.name" type="text" class="form-input" placeholder="如：企业内网拓扑、DMZ区域拓扑" />
          </div>

          <div class="form-section">
            <div class="section-header">
              <h4>节点配置</h4>
              <button type="button" class="btn btn-sm btn-secondary" @click="addNode">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                添加节点
              </button>
            </div>
            <div class="nodes-table-wrapper" v-if="topoForm.nodes.length">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>节点标识</th>
                    <th>显示名称</th>
                    <th>节点类型</th>
                    <th>区域</th>
                    <th>层级</th>
                    <th>关键性</th>
                    <th>绑定产品</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(node, idx) in topoForm.nodes" :key="idx">
                    <td class="mono">
                      <input v-model="node.node_key" type="text" class="inline-input" placeholder="如: firewall-1" />
                    </td>
                    <td>
                      <input v-model="node.name" type="text" class="inline-input" placeholder="如: 主防火墙" />
                    </td>
                    <td>
                      <select v-model="node.node_type" class="inline-select" @change="onNodeTypeChange(node)">
                        <option value="hardware">硬件设备</option>
                        <option value="software">软件系统</option>
                        <option value="os">操作系统</option>
                        <option value="service">服务</option>
                      </select>
                    </td>
                    <td>
                      <select v-model="node.zone" class="inline-select">
                        <option value="internet">Internet</option>
                        <option value="edge">Edge</option>
                        <option value="dmz">DMZ</option>
                        <option value="internal">Internal</option>
                        <option value="core">Core</option>
                        <option value="data">Data</option>
                      </select>
                    </td>
                    <td>
                      <input v-model.number="node.layer" type="number" class="inline-input narrow" min="1" max="10" />
                    </td>
                    <td>
                      <select v-model="node.criticality" class="inline-select">
                        <option value="low">低</option>
                        <option value="normal">普通</option>
                        <option value="high">高</option>
                        <option value="critical">危急</option>
                      </select>
                    </td>
                    <td>
                      <select v-model="node.product_id" class="inline-select" :disabled="node.node_type !== 'hardware'">
                        <option value="">-- 不绑定 --</option>
                        <option v-for="p in products" :key="p.id" :value="p.id">
                          {{ p.name }} ({{ getTypeName(p.type_id) }})
                        </option>
                      </select>
                    </td>
                    <td>
                      <button type="button" class="action-btn delete" @click="removeNode(idx)" title="删除节点">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <line x1="18" y1="6" x2="6" y2="18"/>
                          <line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else class="empty-hint">
              暂无节点，请点击"添加节点"开始配置
            </div>
          </div>

          <div class="form-section">
            <div class="section-header">
              <h4>边配置</h4>
              <button type="button" class="btn btn-sm btn-secondary" @click="addEdge" :disabled="!canAddEdge">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                添加边
              </button>
            </div>
            <div class="edges-table-wrapper" v-if="topoForm.edges.length">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>源节点</th>
                    <th>目标节点</th>
                    <th>边类型</th>
                    <th>方向</th>
                    <th>权重</th>
                    <th>风险值</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(edge, idx) in topoForm.edges" :key="idx">
                    <td>
                      <select v-model="edge.from_node_key" class="inline-select">
                        <option value="">选择源节点</option>
                        <option v-for="n in topoForm.nodes" :key="n.node_key" :value="n.node_key">
                          {{ n.name || n.node_key }}
                        </option>
                      </select>
                    </td>
                    <td>
                      <select v-model="edge.to_node_key" class="inline-select">
                        <option value="">选择目标节点</option>
                        <option v-for="n in topoForm.nodes" :key="n.node_key" :value="n.node_key">
                          {{ n.name || n.node_key }}
                        </option>
                      </select>
                    </td>
                    <td>
                      <select v-model="edge.edge_type" class="inline-select">
                        <option value="network">网络连接</option>
                        <option value="install">部署关系</option>
                        <option value="depend">依赖关系</option>
                        <option value="trust">信任关系</option>
                      </select>
                    </td>
                    <td>
                      <select v-model="edge.direction" class="inline-select">
                        <option value="uni">单向</option>
                        <option value="bi">双向</option>
                      </select>
                    </td>
                    <td>
                      <input v-model.number="edge.weight" type="number" class="inline-input narrow" min="1" max="10" />
                    </td>
                    <td>
                      <input v-model.number="edge.risk" type="number" class="inline-input narrow" min="0" max="100" />
                    </td>
                    <td>
                      <button type="button" class="action-btn delete" @click="removeEdge(idx)" title="删除边">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <line x1="18" y1="6" x2="6" y2="18"/>
                          <line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else class="empty-hint">
              暂无边配置{{ topoForm.nodes.length < 2 ? '（需要至少2个节点才能添加边）' : '（边可自动生成，也可手动添加）' }}
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeTopoModal">取消</button>
          <button type="button" class="btn btn-primary" @click="saveTopo" :disabled="saving">
            {{ saving ? '保存中...' : (editingTopo ? '保存' : '创建') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Graphic Edit Modal (图形化模式) -->
    <div v-if="showGraphicModal" class="modal-overlay" @click.self="closeGraphicModal">
      <div class="modal modal-fullscreen">
        <div class="modal-header">
          <h3>{{ graphicTopoName }} - 图形编辑</h3>
          <div class="graphic-toolbar">
            <button type="button" class="btn btn-sm" :class="editMode === 'select' ? 'btn-primary' : 'btn-secondary'" @click="editMode = 'select'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 3l7.07 16.97 2.51-7.39 7.39-2.51L3 3z"/>
              </svg>
              选择
            </button>
            <button type="button" class="btn btn-sm" :class="editMode === 'add' ? 'btn-primary' : 'btn-secondary'" @click="editMode = 'add'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="16"/>
                <line x1="8" y1="12" x2="16" y2="12"/>
              </svg>
              添加节点
            </button>
            <button type="button" class="btn btn-sm" :class="editMode === 'connect' ? 'btn-primary' : 'btn-secondary'" @click="editMode = 'connect'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                <polyline points="15 3 21 3 21 9"/>
                <line x1="10" y1="14" x2="21" y2="3"/>
              </svg>
              连接
            </button>
            <button type="button" class="btn btn-sm" :class="editMode === 'delete' ? 'btn-primary' : 'btn-secondary'" @click="editMode = 'delete'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
              删除
            </button>
            <div class="toolbar-divider"></div>
            <div class="node-type-selector" v-if="editMode === 'add'">
              <select v-model="newNodeType" class="inline-select">
                <option value="hardware">硬件设备</option>
                <option value="software">软件系统</option>
                <option value="os">操作系统</option>
                <option value="service">服务</option>
              </select>
            </div>
            <div class="edge-type-selector" v-if="editMode === 'connect'">
              <select v-model="newEdgeType" class="inline-select">
                <option value="network">网络连接</option>
                <option value="install">部署关系</option>
                <option value="depend">依赖关系</option>
                <option value="trust">信任关系</option>
              </select>
              <select v-model="newEdgeDirection" class="inline-select">
                <option value="uni">单向</option>
                <option value="bi">双向</option>
              </select>
            </div>
          </div>
          <button class="modal-close" @click="closeGraphicModal">×</button>
        </div>
        <div class="graphic-body">
          <div class="graph-container" ref="graphicContainer"></div>
          <div class="node-properties" v-if="selectedGraphicNode">
            <h4>节点属性</h4>
            <div class="property-form">
              <div class="form-group">
                <label>节点标识</label>
                <input v-model="selectedGraphicNode.node_key" type="text" class="form-input" />
              </div>
              <div class="form-group">
                <label>显示名称</label>
                <input v-model="selectedGraphicNode.name" type="text" class="form-input" />
              </div>
              <div class="form-group">
                <label>节点类型</label>
                <select v-model="selectedGraphicNode.node_type" class="form-input" @change="onGraphicNodeTypeChange">
                  <option value="hardware">硬件设备</option>
                  <option value="software">软件系统</option>
                  <option value="os">操作系统</option>
                  <option value="service">服务</option>
                </select>
              </div>
              <div class="form-group">
                <label>区域</label>
                <select v-model="selectedGraphicNode.zone" class="form-input">
                  <option value="internet">Internet</option>
                  <option value="edge">Edge</option>
                  <option value="dmz">DMZ</option>
                  <option value="internal">Internal</option>
                  <option value="core">Core</option>
                  <option value="data">Data</option>
                </select>
              </div>
              <div class="form-group">
                <label>层级</label>
                <input v-model.number="selectedGraphicNode.layer" type="number" class="form-input" min="1" max="10" />
              </div>
              <div class="form-group">
                <label>关键性</label>
                <select v-model="selectedGraphicNode.criticality" class="form-input">
                  <option value="low">低</option>
                  <option value="normal">普通</option>
                  <option value="high">高</option>
                  <option value="critical">危急</option>
                </select>
              </div>
              <div class="form-group" v-if="selectedGraphicNode.node_type === 'hardware'">
                <label>绑定产品</label>
                <select v-model="selectedGraphicNode.product_id" class="form-input">
                  <option value="">-- 不绑定 --</option>
                  <option v-for="p in products" :key="p.id" :value="p.id">
                    {{ p.name }} ({{ getTypeName(p.type_id) }})
                  </option>
                </select>
              </div>
            </div>
          </div>
          <div class="edge-properties" v-if="selectedGraphicEdge && !selectedGraphicNode">
            <h4>边属性</h4>
            <div class="property-form">
              <div class="form-group">
                <label>源节点</label>
                <input :value="selectedGraphicEdge.from" type="text" class="form-input" disabled />
              </div>
              <div class="form-group">
                <label>目标节点</label>
                <input :value="selectedGraphicEdge.to" type="text" class="form-input" disabled />
              </div>
              <div class="form-group">
                <label>边类型</label>
                <select v-model="selectedGraphicEdge.edgeType" class="form-input">
                  <option value="network">网络连接</option>
                  <option value="install">部署关系</option>
                  <option value="depend">依赖关系</option>
                  <option value="trust">信任关系</option>
                </select>
              </div>
              <div class="form-group">
                <label>方向</label>
                <select v-model="selectedGraphicEdge.direction" class="form-input">
                  <option value="uni">单向</option>
                  <option value="bi">双向</option>
                </select>
              </div>
              <div class="form-group">
                <label>权重</label>
                <input v-model.number="selectedGraphicEdge.weight" type="number" class="form-input" min="1" max="10" />
              </div>
              <div class="form-group">
                <label>风险值</label>
                <input v-model.number="selectedGraphicEdge.risk" type="number" class="form-input" min="0" max="100" />
              </div>
            </div>
          </div>
          <div class="graph-hint" v-if="!selectedGraphicNode && !selectedGraphicEdge">
            <p v-if="editMode === 'select'">点击节点或边进行选择</p>
            <p v-if="editMode === 'add'">点击画布空白处添加节点</p>
            <p v-if="editMode === 'connect'">先点击源节点，再点击目标节点创建连接</p>
            <p v-if="editMode === 'delete'">点击节点或边进行删除</p>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeGraphicModal">取消</button>
          <button type="button" class="btn btn-primary" @click="saveGraphicChanges" :disabled="saving">
            {{ saving ? '保存中...' : '保存修改' }}
          </button>
        </div>
      </div>
    </div>

    <!-- View Modal -->
    <div v-if="showViewModal" class="modal-overlay" @click.self="closeViewModal">
      <div class="modal modal-lg">
        <div class="modal-header">
          <h3>{{ viewTopo?.name }}</h3>
          <button class="modal-close" @click="closeViewModal">×</button>
        </div>
        <div class="modal-body">
          <div class="view-stats">
            <div class="view-stat">
              <span class="stat-label">节点数</span>
              <span class="stat-value">{{ viewTopo?.nodes?.length || 0 }}</span>
            </div>
            <div class="view-stat">
              <span class="stat-label">边数</span>
              <span class="stat-value">{{ viewTopo?.edges?.length || 0 }}</span>
            </div>
            <div class="view-stat">
              <span class="stat-label">设备数</span>
              <span class="stat-value">{{ viewTopo?.products?.length || 0 }}</span>
            </div>
          </div>

          <div class="view-section" v-if="viewTopo?.nodes?.length">
            <h4>节点列表</h4>
            <div class="view-nodes">
              <div v-for="node in viewTopo.nodes" :key="node.node_key" class="view-node-card">
                <div class="node-card-header">
                  <span class="node-type-badge" :class="node.node_type">{{ node.node_type }}</span>
                  <span class="node-zone">{{ node.zone }}</span>
                </div>
                <div class="node-card-name">{{ node.name }}</div>
                <div class="node-card-key">{{ node.node_key }}</div>
                <div class="node-card-meta">
                  <span>层级: {{ node.layer }}</span>
                  <span>关键性: {{ node.criticality }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="view-section" v-if="viewTopo?.edges?.length">
            <h4>边列表</h4>
            <div class="view-edges">
              <div v-for="edge in viewTopo.edges" :key="`${edge.from_node_key}-${edge.to_node_key}`" class="view-edge-item">
                <span class="edge-from">{{ edge.from_node_key }}</span>
                <svg class="edge-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M5 12h14M12 5l7 7-7 7"/>
                </svg>
                <span class="edge-to">{{ edge.to_node_key }}</span>
                <span class="edge-type">{{ edge.edge_type }}</span>
                <span class="edge-direction">{{ edge.direction === 'bi' ? '双向' : '单向' }}</span>
              </div>
            </div>
          </div>
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
          <p class="delete-message">确定要删除拓扑 "{{ deleteTarget?.name }}" 吗？此操作将同时删除所有关联的节点和边，无法撤销。</p>
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
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { topoApi, productApi } from '../api'
import { Network, DataSet } from 'vis-network/standalone'
import 'vis-network/styles/vis-network.css'
import { useTheme } from '../composables/useTheme'

const { isDark } = useTheme()

const nodeFontColor = computed(() => isDark.value ? '#e6edf3' : '#1a1f2e')

const searchQuery = ref('')
const topologies = ref([])
const products = ref([])

const showTopoModal = ref(false)
const editingTopo = ref(null)
const saving = ref(false)
const showViewModal = ref(false)
const viewTopo = ref(null)
const showDeleteModal = ref(false)
const deleteTarget = ref(null)
const showGraphicModal = ref(false)
const graphicTopoId = ref(null)
const graphicTopoName = ref('')
const editMode = ref('select')
const newNodeType = ref('hardware')
const newEdgeType = ref('network')
const newEdgeDirection = ref('uni')
const selectedGraphicNode = ref(null)
const selectedGraphicEdge = ref(null)
const connectSource = ref(null)
const graphicContainer = ref(null)
let graphicNetwork = null
const graphicNodes = ref([])
const graphicEdges = ref([])

const topoForm = ref({
  name: '',
  nodes: [],
  edges: []
})

const productTypes = ref([])

const filteredTopologies = computed(() => {
  if (!searchQuery.value) return topologies.value
  const q = searchQuery.value.toLowerCase()
  return topologies.value.filter(t => t.name.toLowerCase().includes(q))
})

const canAddEdge = computed(() => topoForm.value.nodes.length >= 2)

function getTypeName(typeId) {
  const t = productTypes.value.find(pt => pt.id === typeId)
  return t ? t.name : '-'
}

function createEmptyNode() {
  return {
    node_key: '',
    name: '',
    node_type: 'hardware',
    zone: 'internal',
    layer: 1,
    criticality: 'normal',
    product_id: ''
  }
}

function createEmptyEdge() {
  return {
    from_node_key: '',
    to_node_key: '',
    edge_type: 'network',
    direction: 'uni',
    weight: 1,
    risk: 0
  }
}

function addNode() {
  topoForm.value.nodes.push(createEmptyNode())
}

function removeNode(idx) {
  const removedKey = topoForm.value.nodes[idx].node_key
  topoForm.value.nodes.splice(idx, 1)
  topoForm.value.edges = topoForm.value.edges.filter(
    e => e.from_node_key !== removedKey && e.to_node_key !== removedKey
  )
}

function addEdge() {
  if (topoForm.value.nodes.length >= 2) {
    topoForm.value.edges.push(createEmptyEdge())
  }
}

function removeEdge(idx) {
  topoForm.value.edges.splice(idx, 1)
}

function onNodeTypeChange(node) {
  if (node.node_type !== 'hardware') {
    node.product_id = ''
  }
}

async function loadData() {
  try {
    const [topoRes, prodRes, typeRes] = await Promise.all([
      topoApi.list({ mode: 'summary' }),
      productApi.listProducts(),
      productApi.list()
    ])
    topologies.value = topoRes.items || topoRes || []
    products.value = prodRes || []
    productTypes.value = typeRes || []
  } catch (err) {
    console.error('Failed to load data:', err)
  }
}

async function loadTopoDetail(id) {
  try {
    return await topoApi.get(id)
  } catch (err) {
    console.error('Failed to load topology detail:', err)
    return null
  }
}

function openTopoModal(item = null) {
  editingTopo.value = item
  if (item) {
    topoForm.value = { name: item.name, nodes: [], edges: [] }
    loadTopoDetail(item.id).then(detail => {
      if (detail) {
        topoForm.value.nodes = (detail.nodes || []).map(n => ({
          node_key: n.node_key,
          name: n.name,
          node_type: n.node_type,
          zone: n.zone,
          layer: n.layer,
          criticality: n.criticality,
          product_id: n.product_id || ''
        }))
        topoForm.value.edges = (detail.edges || []).map(e => ({
          from_node_key: e.from_node_key,
          to_node_key: e.to_node_key,
          edge_type: e.edge_type,
          direction: e.direction,
          weight: e.weight,
          risk: e.risk
        }))
      }
    })
  } else {
    topoForm.value = { name: '', nodes: [], edges: [] }
  }
  showTopoModal.value = true
}

function closeTopoModal() {
  showTopoModal.value = false
  editingTopo.value = null
}

async function saveTopo() {
  if (!topoForm.value.name) {
    alert('请输入拓扑名称')
    return
  }

  const validNodes = topoForm.value.nodes.filter(n => n.node_key && n.name)
  if (validNodes.length === 0) {
    alert('请至少添加一个有效节点（节点标识和名称不能为空）')
    return
  }

  const hasHardware = validNodes.some(n => n.node_type === 'hardware')
  if (!hasHardware) {
    alert('拓扑至少需要一个硬件类型节点')
    return
  }

  saving.value = true
  try {
    const payload = {
      name: topoForm.value.name,
      nodes: validNodes.map(n => ({
        node_key: n.node_key,
        name: n.name,
        node_type: n.node_type,
        zone: n.zone,
        layer: n.layer,
        criticality: n.criticality,
        product_id: n.node_type === 'hardware' && n.product_id ? Number(n.product_id) : null
      })),
      edges: topoForm.value.edges.filter(e => e.from_node_key && e.to_node_key).map(e => ({
        from_node_key: e.from_node_key,
        to_node_key: e.to_node_key,
        edge_type: e.edge_type,
        direction: e.direction,
        weight: e.weight,
        risk: e.risk
      }))
    }

    if (editingTopo.value) {
      await topoApi.update(editingTopo.value.id, payload)
    } else {
      await topoApi.create(payload)
    }

    await loadData()
    closeTopoModal()
  } catch (err) {
    alert(err.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function openGraphicModal(item) {
  graphicTopoId.value = item.id
  graphicTopoName.value = item.name
  editMode.value = 'select'
  selectedGraphicNode.value = null
  selectedGraphicEdge.value = null
  connectSource.value = null

  const detail = await loadTopoDetail(item.id)
  if (!detail) return

  graphicNodes.value = (detail.nodes || []).map(n => ({
    id: n.node_key,
    node_key: n.node_key,
    label: n.name,
    name: n.name,
    node_type: n.node_type,
    zone: n.zone,
    layer: n.layer,
    criticality: n.criticality,
    product_id: n.product_id || '',
    color: getNodeColor(n.node_type)
  }))

  graphicEdges.value = (detail.edges || []).map((e, idx) => ({
    id: `edge-${idx}`,
    from: e.from_node_key,
    to: e.to_node_key,
    edgeType: e.edge_type,
    direction: e.direction,
    weight: e.weight,
    risk: e.risk,
    color: '#4a5568'
  }))

  showGraphicModal.value = true
  await nextTick()
  initGraphicGraph()
}

function getNodeColor(nodeType) {
  const colors = {
    hardware: '#ff6b6b',
    software: '#58a6ff',
    os: '#a371f7',
    service: '#f0883e'
  }
  return colors[nodeType] || '#8b949e'
}

async function initGraphicGraph() {
  if (!graphicContainer.value) return

  if (graphicNetwork) {
    graphicNetwork.destroy()
    graphicNetwork = null
  }

  const nodesData = new DataSet(graphicNodes.value.map(n => ({
    id: n.id,
    label: n.label,
    color: {
      background: n.color + '30',
      border: n.color,
      highlight: { background: n.color + '50', border: n.color }
    },
    font: { color: nodeFontColor.value, size: 14 },
    shape: 'box',
    size: 30
  })))

  const edgesData = new DataSet(graphicEdges.value.map(e => ({
    id: e.id,
    from: e.from,
    to: e.to,
    color: { color: e.color, highlight: '#00d4aa' },
    arrows: e.direction === 'bi' ? 'to;from' : 'to',
    width: e.weight || 1
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

  graphicNetwork = new Network(graphicContainer.value, { nodes: nodesData, edges: edgesData }, options)

  graphicNetwork.on('click', (params) => {
    if (params.nodes.length > 0) {
      const nodeId = params.nodes[0]
      if (editMode.value === 'delete') {
        deleteGraphicNode(nodeId)
      } else if (editMode.value === 'connect') {
        handleConnect(nodeId)
      } else if (editMode.value === 'select') {
        const node = graphicNodes.value.find(n => n.id === nodeId)
        selectedGraphicNode.value = node ? { ...node } : null
        selectedGraphicEdge.value = null
      }
    } else if (params.edges.length > 0) {
      const edgeId = params.edges[0]
      if (editMode.value === 'delete') {
        deleteGraphicEdge(edgeId)
      } else if (editMode.value === 'select') {
        const edge = graphicEdges.value.find(e => e.id === edgeId)
        selectedGraphicEdge.value = edge ? { ...edge } : null
        selectedGraphicNode.value = null
      }
    } else {
      if (editMode.value === 'add') {
        addGraphicNode(params.pointer.canvas)
      } else if (editMode.value === 'select') {
        selectedGraphicNode.value = null
        selectedGraphicEdge.value = null
      }
    }
  })
}

function handleConnect(nodeId) {
  if (!connectSource.value) {
    connectSource.value = nodeId
  } else {
    if (connectSource.value !== nodeId) {
      addGraphicEdge(connectSource.value, nodeId)
    }
    connectSource.value = null
  }
}

function addGraphicNode(position) {
  // 如果选择硬件类型但没有可用的产品，提示用户
  if (newNodeType.value === 'hardware' && products.value.length === 0) {
    alert('当前没有可用的安全产品，请先在"数据管理"中添加产品，或选择其他节点类型（软件系统/操作系统/服务）')
    return
  }

  const nodeCount = graphicNodes.value.length + 1
  const typeLabel = { hardware: '设备', software: '软件', os: '系统', service: '服务' }
  const nodeKey = `node-${nodeCount}`
  const newNode = {
    id: nodeKey,
    node_key: nodeKey,
    label: `${typeLabel[newNodeType.value] || '节点'}${nodeCount}`,
    name: `${typeLabel[newNodeType.value] || '节点'}${nodeCount}`,
    node_type: newNodeType.value,
    zone: 'internal',
    layer: Math.ceil(nodeCount / 3),
    criticality: 'normal',
    product_id: newNodeType.value === 'hardware' && products.value.length > 0 ? products.value[0].id : '',
    color: getNodeColor(newNodeType.value)
  }
  graphicNodes.value.push(newNode)

  const dataset = graphicNetwork.body.data.nodes
  dataset.add({
    id: newNode.id,
    label: newNode.label,
    color: {
      background: newNode.color + '30',
      border: newNode.color,
      highlight: { background: newNode.color + '50', border: newNode.color }
    },
    font: { color: nodeFontColor.value, size: 14 },
    shape: 'box',
    size: 30
  })
}

function addGraphicEdge(fromId, toId) {
  const edgeId = `edge-${Date.now()}`
  const newEdge = {
    id: edgeId,
    from: fromId,
    to: toId,
    edgeType: newEdgeType.value,
    direction: newEdgeDirection.value,
    weight: 1,
    risk: 0,
    color: '#4a5568'
  }
  graphicEdges.value.push(newEdge)

  const dataset = graphicNetwork.body.data.edges
  dataset.add({
    id: edgeId,
    from: fromId,
    to: toId,
    color: { color: '#4a5568', highlight: '#00d4aa' },
    arrows: newEdgeDirection.value === 'bi' ? 'to;from' : 'to',
    width: 1
  })
}

function deleteGraphicNode(nodeId) {
  graphicNodes.value = graphicNodes.value.filter(n => n.id !== nodeId)
  graphicEdges.value = graphicEdges.value.filter(e => e.from !== nodeId && e.to !== nodeId)

  graphicNetwork.body.data.nodes.remove(nodeId)
  graphicNetwork.body.data.edges.forEach(edge => {
    if (edge.from === nodeId || edge.to === nodeId) {
      graphicNetwork.body.data.edges.remove(edge.id)
    }
  })

  selectedGraphicNode.value = null
}

function deleteGraphicEdge(edgeId) {
  graphicEdges.value = graphicEdges.value.filter(e => e.id !== edgeId)
  graphicNetwork.body.data.edges.remove(edgeId)
  selectedGraphicEdge.value = null
}

function onGraphicNodeTypeChange() {
  if (selectedGraphicNode.value) {
    selectedGraphicNode.value.color = getNodeColor(selectedGraphicNode.value.node_type)
    selectedGraphicNode.value.product_id = ''

    const nodeData = graphicNetwork.body.data.nodes.get(selectedGraphicNode.value.id)
    if (nodeData) {
      nodeData.color = {
        background: selectedGraphicNode.value.color + '30',
        border: selectedGraphicNode.value.color
      }
      graphicNetwork.body.data.nodes.update(nodeData)
    }
  }
}

async function saveGraphicChanges() {
  if (!graphicTopoId.value) return

  saving.value = true
  try {
    const validNodes = graphicNodes.value.filter(n => n.node_key && n.name)
    const payload = {
      name: graphicTopoName.value,
      nodes: validNodes.map(n => ({
        node_key: n.node_key,
        name: n.name,
        node_type: n.node_type,
        zone: n.zone,
        layer: n.layer,
        criticality: n.criticality,
        product_id: n.node_type === 'hardware' && n.product_id ? Number(n.product_id) : null
      })),
      edges: graphicEdges.value.map(e => ({
        from_node_key: e.from,
        to_node_key: e.to,
        edge_type: e.edgeType,
        direction: e.direction,
        weight: e.weight,
        risk: e.risk
      }))
    }

    await topoApi.update(graphicTopoId.value, payload)
    await loadData()
    closeGraphicModal()
  } catch (err) {
    alert(err.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

function closeGraphicModal() {
  showGraphicModal.value = false
  if (graphicNetwork) {
    graphicNetwork.destroy()
    graphicNetwork = null
  }
  graphicTopoId.value = null
  graphicNodes.value = []
  graphicEdges.value = []
  selectedGraphicNode.value = null
  selectedGraphicEdge.value = null
  connectSource.value = null
}

async function openViewModal(item) {
  const detail = await loadTopoDetail(item.id)
  viewTopo.value = detail
  showViewModal.value = true
}

function closeViewModal() {
  showViewModal.value = false
  viewTopo.value = null
}

async function copyTopo(item) {
  try {
    await topoApi.copy(item.id)
    await loadData()
  } catch (err) {
    alert(err.response?.data?.message || '复制失败')
  }
}

function deleteTopo(item) {
  deleteTarget.value = item
  showDeleteModal.value = true
}

function closeDeleteModal() {
  showDeleteModal.value = false
  deleteTarget.value = null
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    await topoApi.delete(deleteTarget.value.id)
    await loadData()
    closeDeleteModal()
  } catch (err) {
    alert(err.response?.data?.message || '删除失败')
  }
}

onMounted(loadData)

watch(isDark, () => {
  if (graphicNetwork && graphicNodes.value.length) {
    nextTick(() => {
      const fontColor = isDark.value ? '#e6edf3' : '#1a1f2e'
      const updatedNodes = []
      graphicNetwork.body.data.nodes.forEach(node => {
        updatedNodes.push({ id: node.id, font: { color: fontColor, size: 14 } })
      })
      graphicNetwork.body.data.nodes.update(updatedNodes)
    })
  }
})

onUnmounted(() => {
  if (graphicNetwork) {
    graphicNetwork.destroy()
    graphicNetwork = null
  }
})
</script>

<style scoped>
.manage-page {
  min-height: 100vh;
}

.page-header {
  padding: 32px 40px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-content {
  padding: 32px 40px;
}

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

.stat-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-family: var(--font-mono);
}

.stat-badge.nodes {
  background: rgba(0, 212, 170, 0.1);
  color: var(--accent-primary);
}

.stat-badge.edges {
  background: rgba(88, 166, 255, 0.1);
  color: var(--info);
}

.stat-badge.products {
  background: rgba(163, 113, 247, 0.1);
  color: #a371f7;
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

.action-btn.view:hover {
  border-color: var(--info);
  color: var(--info);
}

.action-btn.graphic:hover {
  border-color: var(--accent-primary);
  color: var(--accent-primary);
}

.action-btn.copy:hover {
  border-color: #a371f7;
  color: #a371f7;
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

.modal-fullscreen {
  width: 95%;
  height: 90vh;
  max-width: 1400px;
  display: flex;
  flex-direction: column;
}

.modal-lg {
  max-width: 800px;
}

.modal-xl {
  max-width: 1100px;
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
  flex: 1;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
}

/* Graphic Edit */
.graphic-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  justify-content: center;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: var(--border-color);
  margin: 0 8px;
}

.node-type-selector,
.edge-type-selector {
  display: flex;
  gap: 8px;
}

.graphic-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
}

.graph-container {
  flex: 1;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  margin: 16px;
}

.node-properties,
.edge-properties {
  width: 280px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  margin: 16px 16px 16px 0;
  padding: 20px;
  overflow-y: auto;
}

.node-properties h4,
.edge-properties h4 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.property-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.graph-hint {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  padding: 12px 24px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 14px;
  box-shadow: var(--shadow-md);
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

.form-section {
  margin-bottom: 24px;
  padding: 20px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.nodes-table-wrapper,
.edges-table-wrapper {
  overflow-x: auto;
}

.nodes-table-wrapper table,
.edges-table-wrapper table {
  min-width: 900px;
}

.inline-input,
.inline-select {
  padding: 6px 8px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  min-width: 80px;
}

.inline-input:focus,
.inline-select:focus {
  outline: none;
  border-color: var(--accent-primary);
}

.inline-input.narrow {
  width: 60px;
  min-width: 60px;
  text-align: center;
}

.inline-select {
  min-width: 100px;
}

.empty-hint {
  text-align: center;
  padding: 30px;
  color: var(--text-tertiary);
  font-size: 14px;
}

/* View Modal */
.view-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
}

.view-stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.view-stat .stat-label {
  font-size: 12px;
  color: var(--text-tertiary);
}

.view-stat .stat-value {
  font-size: 24px;
  font-weight: 700;
  font-family: var(--font-mono);
  color: var(--accent-primary);
}

.view-section {
  margin-bottom: 20px;
}

.view-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.view-nodes {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.view-node-card {
  padding: 14px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.node-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.node-type-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
}

.node-type-badge.hardware {
  background: rgba(255, 107, 107, 0.1);
  color: #ff6b6b;
}

.node-type-badge.software {
  background: rgba(88, 166, 255, 0.1);
  color: #58a6ff;
}

.node-type-badge.os {
  background: rgba(163, 113, 247, 0.1);
  color: #a371f7;
}

.node-type-badge.service {
  background: rgba(240, 136, 62, 0.1);
  color: #f0883e;
}

.node-zone {
  font-size: 11px;
  color: var(--text-tertiary);
  text-transform: uppercase;
}

.node-card-name {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.node-card-key {
  font-size: 12px;
  font-family: var(--font-mono);
  color: var(--text-tertiary);
  margin-bottom: 8px;
}

.node-card-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-secondary);
}

.view-edges {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.view-edge-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 13px;
}

.edge-from,
.edge-to {
  font-family: var(--font-mono);
  color: var(--text-primary);
}

.edge-arrow {
  width: 16px;
  height: 16px;
  color: var(--text-tertiary);
}

.edge-type,
.edge-direction {
  padding: 2px 8px;
  background: var(--bg-card);
  border-radius: 4px;
  font-size: 11px;
  color: var(--text-secondary);
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

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
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

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

.btn-sm svg {
  width: 14px;
  height: 14px;
}
</style>
