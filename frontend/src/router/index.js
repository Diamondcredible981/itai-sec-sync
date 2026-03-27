import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/topology',
      name: 'topology',
      component: () => import('../views/TopologyView.vue')
    },
    {
      path: '/topology/:id',
      name: 'topology-detail',
      component: () => import('../views/TopologyView.vue')
    },
    {
      path: '/analysis',
      name: 'analysis',
      component: () => import('../views/AnalysisView.vue')
    },
    {
      path: '/suggest/:id',
      name: 'suggest',
      component: () => import('../views/SuggestView.vue')
    }
  ]
})

export default router
