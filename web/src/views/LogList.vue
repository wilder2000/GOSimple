<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="mb-0">登录日志</h4>
      <small class="text-muted">共 {{ totalRows }} 条</small>
    </div>
    <div class="card card-table">
      <div class="table-responsive">
        <table class="table table-hover mb-0">
          <thead class="table-light"><tr><th>ID</th><th>账号</th><th>IP</th><th>登录时间</th></tr></thead>
          <tbody>
            <tr v-for="l in list" :key="l.id">
              <td>{{ l.id }}</td><td>{{ l.account }}</td><td>{{ l.ip }}</td><td>{{ l.logintime?.slice(0, 19) }}</td>
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
import { query } from '../api/mif'
import Pagination from '../components/Pagination.vue'

const list = ref<any[]>([])
const page = ref(1)
const totalPages = ref(0)
const totalRows = ref(0)

async function loadData(p?: number) {
  if (p) page.value = p
  const res = await query({ Target: 'log', PageIndex: page.value, PageSize: 15, Order: 'id desc' })
  list.value = res.data.Data || []
  totalPages.value = res.data.TotalPages
  totalRows.value = res.data.TotalRows
}

onMounted(() => loadData())
</script>
