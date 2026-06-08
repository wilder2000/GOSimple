<template>
  <div class="d-flex">
    <aside class="sidebar d-flex flex-column py-3">
      <div class="text-center mb-4">
        <h5 class="text-white mb-0">GOSimple</h5>
        <small class="text-secondary">后台管理</small>
      </div>
      <nav class="nav flex-column">
        <router-link v-for="item in menu" :key="item.path" :to="item.path" class="nav-link" :class="{ active: isActive(item.path) }">
          <i :class="item.icon"></i>{{ item.label }}
        </router-link>
      </nav>
      <div class="mt-auto px-3">
        <hr class="border-secondary">
        <a href="#" class="nav-link" @click.prevent="logout">
          <i class="bi bi-box-arrow-right"></i>退出登录
        </a>
      </div>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const menu = [
  { path: '/dashboard', label: '仪表盘', icon: 'bi bi-speedometer2' },
  { path: '/users', label: '用户管理', icon: 'bi bi-people' },
  { path: '/roles', label: '角色管理', icon: 'bi bi-shield' },
  { path: '/groups', label: '编组管理', icon: 'bi bi-layers' },
  { path: '/departments', label: '部门管理', icon: 'bi bi-building' },
  { path: '/operators', label: '功能权限', icon: 'bi bi-key' },
  { path: '/urlmappings', label: 'URL 映射', icon: 'bi bi-link-45deg' },
  { path: '/logs', label: '登录日志', icon: 'bi bi-clock-history' },
]

function isActive(path: string) {
  return route.path.startsWith(path)
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>
