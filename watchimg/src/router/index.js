import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '@/components/Index/Index'
import AppList from '@/components/Account/AccountList'
import AppListd from '@/components/Account/AccdList'
import BanAccdList from '@/components/Account/BanAccdList'
import AccdDateList from '@/components/Account/AccdDateList'

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
  {
    path: '/banacclist',
    name: 'banacclist',
    component: BanAccdList
  },
  {
    path: '/accdatelist',
    name: 'accdatelist',
    component: AccdDateList
  },
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router