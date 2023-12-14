import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import GraphView from '../views/GraphView.vue'
import HistoryView from '@/views/HistoryView.vue'
import AnalyticsView from '@/views/AnalyticsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/graph',
      name: 'graph',
      component: GraphView
    },
  {
    path: '/history',
    name: 'history',
    component: HistoryView
  },
  {
    path: '/analytics',
    name: 'analytics',
    component: AnalyticsView
  },
  ]
})

export default router
