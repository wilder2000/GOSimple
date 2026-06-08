<template>
  <div>
    <h4 class="mb-4">仪表盘</h4>
    <div class="row g-3">
      <div v-for="card in cards" :key="card.label" class="col-md-3">
        <div class="card border-0 shadow-sm">
          <div class="card-body text-center py-4">
            <i :class="card.icon" style="font-size: 2rem; color: var(--bs-primary)"></i>
            <h3 class="mt-2 mb-0">{{ card.count }}</h3>
            <small class="text-muted">{{ card.label }}</small>
          </div>
        </div>
      </div>
    </div>
    <div class="card card-table mt-4">
      <div class="card-header">
        <h6 class="mb-0">快速入口</h6>
      </div>
      <div class="card-body">
        <div class="row g-2">
          <div class="col-md-3">
            <router-link to="/users" class="btn btn-outline-primary w-100">
              <i class="bi bi-people me-1"></i>用户管理
            </router-link>
          </div>
          <div class="col-md-3">
            <router-link to="/roles" class="btn btn-outline-success w-100">
              <i class="bi bi-shield me-1"></i>角色管理
            </router-link>
          </div>
          <div class="col-md-3">
            <router-link to="/groups" class="btn btn-outline-info w-100">
              <i class="bi bi-layers me-1"></i>编组管理
            </router-link>
          </div>
          <div class="col-md-3">
            <router-link to="/logs" class="btn btn-outline-secondary w-100">
              <i class="bi bi-clock-history me-1"></i>登录日志
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { query } from '../api/mif'

const cards = ref([
  { label: '用户总数', icon: 'bi bi-people', count: 0 },
  { label: '角色数', icon: 'bi bi-shield', count: 0 },
  { label: '编组数', icon: 'bi bi-layers', count: 0 },
  { label: '部门数', icon: 'bi bi-building', count: 0 },
])

onMounted(async () => {
  try {
    const [users, roles, groups, depts] = await Promise.all([
      query({ Target: 'user', PageSize: 1 }),
      query({ Target: 'role', PageSize: 1 }),
      query({ Target: 'group', PageSize: 1 }),
      query({ Target: 'depart', PageSize: 1 }),
    ])
    cards.value[0].count = users.data.TotalRows
    cards.value[1].count = roles.data.TotalRows
    cards.value[2].count = groups.data.TotalRows
    cards.value[3].count = depts.data.TotalRows
  } catch (_) { /* ignore */ }
})
</script>
