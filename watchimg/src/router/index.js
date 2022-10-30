import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '@/components/Index/Index'
import AppList from '@/components/Account/AccountList'
import AppListd from '@/components/Account/AccdList'

const routerHistory = createWebHistory()

const constantRoutes = [
  {
    path: '/',
    name: 'index',
    component: AppIndex
  },
  {
    path: '/acclist',
    name: 'acclist',
    component: AppList
  },
  {
    path: '/acclistd',
    name: 'acclistd',
    component: AppListd
  },
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router