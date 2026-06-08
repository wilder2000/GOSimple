<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="mb-0">用户管理</h4>
      <router-link to="/users/new" class="btn btn-primary btn-sm">
        <i class="bi bi-plus-lg"></i> 新增用户
      </router-link>
    </div>
    <div class="card card-table">
      <div class="card-header d-flex justify-content-between align-items-center">
        <div class="d-flex gap-2">
          <input v-model="search" class="form-control form-control-sm" style="width:200px" placeholder="搜索邮箱/姓名" @keyup.enter="loadData(1)">
          <button class="btn btn-sm btn-outline-secondary" @click="loadData(1)"><i class="bi bi-search"></i></button>
        </div>
        <small class="text-muted">共 {{ totalRows }} 条</small>
      </div>
      <div class="table-responsive">
        <table class="table table-hover mb-0">
          <thead class="table-light">
            <tr>
              <th>邮箱</th><th>姓名</th><th>手机</th><th>性别</th><th>状态</th><th>创建时间</th><th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in list" :key="u.id">
              <td>{{ u.email }}</td>
              <td>{{ u.name || '-' }}</td>
              <td>{{ u.mobile || '-' }}</td>
              <td>{{ u.sex === 1 ? '男' : u.sex === 0 ? '女' : '未知' }}</td>
              <td><span class="badge" :class="u.state === 1 ? 'bg-success' : 'bg-secondary'">{{ u.state === 1 ? '启用' : '禁用' }}</span></td>
              <td>{{ u.createtime?.slice(0, 10) }}</td>
              <td>
                <router-link :to="`/users/${u.id}/edit`" class="btn btn-sm btn-outline-primary me-1"><i class="bi bi-pencil"></i></router-link>
                <button class="btn btn-sm btn-outline-danger" @click="del(u)"><i class="bi bi-trash"></i></button>
              </td>
            </tr>
            <tr v-if="!list.length"><td colspan="7" class="text-center text-muted py-4">暂无数据</td></tr>
          </tbody>
        </table>
      </div>
      <div class="card-footer d-flex justify-content-center">
        <Pagination v-model="page" :total-pages="totalPages" @update:model-value="loadData" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { query, remove } from '../api/mif'
import Pagination from '../components/Pagination.vue'

const list = ref<any[]>([])
const page = ref(1)
const totalPages = ref(0)
const totalRows = ref(0)
const search = ref('')
const pageSize = 15

async function loadData(p?: number) {
  if (p) page.value = p
  const where: Record<string, any> = {}
  if (search.value) {
    where['name like'] = `%${search.value}%`
  }
  const res = await query({
    Target: 'user',
    Where: where,
    PageIndex: page.value,
    PageSize: pageSize,
  })
  list.value = res.data.Data || []
  totalPages.value = res.data.TotalPages
  totalRows.value = res.data.TotalRows
}

async function del(u: any) {
  if (!confirm(`确定删除用户 ${u.email}？`)) return
  await remove('user', { id: u.id })
  loadData()
}

onMounted(() => loadData())
</script>
