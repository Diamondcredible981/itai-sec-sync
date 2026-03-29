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
      path: '/suggest',
      name: 'suggest',
      component: () => import('../views/SuggestView.vue')
    },
    {
      path: '/suggest/:id',
      name: 'suggest-detail',
      component: () => import('../views/SuggestView.vue')
    },
    {
      path: '/manage',
      name: 'manage',
      component: () => import('../views/ManageView.vue')
    },
    {
      path: '/topology-manage',
      name: 'topology-manage',
      component: () => import('../views/TopologyManageView.vue')
    }
  ]
})

export default router
