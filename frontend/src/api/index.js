import axios from 'axios'

const api = axios.create({
  baseURL: '/',
  timeout: 10000
})

api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// 产品类型 CRUD
export const productTypeApi = {
  list: () => api.get('/type'),
  get: (id) => api.get(`/type/${id}`),
  create: (data) => api.post('/type', data),
  update: (id, data) => api.put(`/type/${id}`, data),
  delete: (id) => api.delete(`/type/${id}`)
}

// 功能点 CRUD
export const functionApi = {
  list: () => api.get('/function'),
  get: (id) => api.get(`/function/${id}`),
  create: (data) => api.post('/function', data),
  update: (id, data) => api.put(`/function/${id}`, data),
  delete: (id) => api.delete(`/function/${id}`),
  getByCategory: (category) => api.get('/functions/by-category', { params: { category } })
}

// 产品 CRUD
export const productApi = {
  list: (params) => api.get('/type', { params }),
  listProducts: (params) => api.get('/product', { params }),
  getProduct: (id) => api.get(`/product/${id}`),
  create: (data) => api.post('/product', data),
  update: (id, data) => api.put(`/product/${id}`, data),
  delete: (id) => api.delete(`/product/${id}`),
  getProductsByIds: (ids) => api.post('/products/batch', { product_ids: ids })
}

export const topoApi = {
  list: (params) => api.get('/topo', { params }),
  get: (id) => api.get(`/topo/${id}`),
  create: (data) => api.post('/topo', data),
  update: (id, data) => api.put(`/topo/${id}`, data),
  delete: (id) => api.delete(`/topo/${id}`),
  copy: (id) => api.post(`/topo/${id}/copy`),
  visualization: (id) => api.get(`/topo/${id}/visualization`),
  attackPath: (id, params) => api.get(`/topo/${id}/attack-path`, { params })
}

export const analyzeApi = {
  byProductIds: (productIds) => api.post('/analyze', { product_ids: productIds }),
  byTopoId: (id) => api.post(`/analyze/by-topo/${id}`)
}

export const suggestApi = {
  get: (id, strategy = 'min-change') => api.post(`/suggest/${id}?strategy=${strategy}`)
}

export default api
