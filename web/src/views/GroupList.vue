<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="mb-0">编组管理</h4>
      <router-link to="/groups/new" class="btn btn-primary btn-sm"><i class="bi bi-plus-lg"></i> 新增编组</router-link>
    </div>
    <div class="card card-table">
      <div class="card-header d-flex justify-content-between"><small class="text-muted">共 {{ totalRows }} 条</small></div>
      <div class="table-responsive">
        <table class="table table-hover mb-0">
          <thead class="table-light"><tr><th>ID</th><th>名称</th><th>创建时间</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="g in list" :key="g.id">
              <td>{{ g.id }}</td><td>{{ g.name }}</td><td>{{ g.createtime?.slice(0, 10) }}</td>
              <td>
                <router-link :to="`/groups/${g.id}/edit`" class="btn btn-sm btn-outline-primary me-1"><i class="bi bi-pencil"></i></router-link>
                <button class="btn btn-sm btn-outline-danger" @click="del(g)"><i class="bi bi-trash"></i></button>
              </td>
            </tr>
            <tr v-if="!list.length"><td colspan="4" class="text-center text-muted py-4">暂无数据</td></tr>
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

async function loadData(p?: number) {
  if (p) page.value = p
  const res = await query({ Target: 'group', PageIndex: page.value, PageSize: 15 })
  list.value = res.data.Data || []
  totalPages.value = res.data.TotalPages
  totalRows.value = res.data.TotalRows
}

async function del(g: any) {
  if (!confirm(`确定删除编组 ${g.name}？`)) return
  await remove('group', { id: g.id })
  loadData()
}
onMounted(() => loadData())
</script>
