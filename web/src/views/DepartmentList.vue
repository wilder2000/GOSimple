<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="mb-0">部门管理</h4>
      <button class="btn btn-primary btn-sm" @click="showForm = true; form = { name: '' }"><i class="bi bi-plus-lg"></i> 新增部门</button>
    </div>
    <div class="card card-table">
      <div class="card-header d-flex justify-content-between"><small class="text-muted">共 {{ totalRows }} 条</small></div>
      <div class="table-responsive">
        <table class="table table-hover mb-0">
          <thead class="table-light"><tr><th>ID</th><th>名称</th><th>创建时间</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="d in list" :key="d.id">
              <td>{{ d.id }}</td><td>{{ d.name }}</td><td>{{ d.createtime?.slice(0, 10) }}</td>
              <td>
                <button class="btn btn-sm btn-outline-primary me-1" @click="edit(d)"><i class="bi bi-pencil"></i></button>
                <button class="btn btn-sm btn-outline-danger" @click="del(d)"><i class="bi bi-trash"></i></button>
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

    <div v-if="showForm" class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.4)">
      <div class="modal-dialog modal-sm modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header"><h6 class="mb-0">{{ editingId ? '编辑部门' : '新增部门' }}</h6><button type="button" class="btn-close" @click="showForm = false"></button></div>
          <div class="modal-body">
            <input v-model="form.name" class="form-control" placeholder="部门名称" @keyup.enter="saveDept">
          </div>
          <div class="modal-footer">
            <button class="btn btn-sm btn-secondary" @click="showForm = false">取消</button>
            <button class="btn btn-sm btn-primary" @click="saveDept">保存</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { query, create, update, remove } from '../api/mif'
import Pagination from '../components/Pagination.vue'

const list = ref<any[]>([])
const page = ref(1)
const totalPages = ref(0)
const totalRows = ref(0)
const showForm = ref(false)
const form = ref<any>({ name: '' })
const editingId = ref<number | null>(null)

async function loadData(p?: number) {
  if (p) page.value = p
  const res = await query({ Target: 'depart', PageIndex: page.value, PageSize: 15 })
  list.value = res.data.Data || []
  totalPages.value = res.data.TotalPages
  totalRows.value = res.data.TotalRows
}

function edit(d: any) {
  editingId.value = d.id
  form.value = { name: d.name }
  showForm.value = true
}

async function saveDept() {
  if (!form.value.name) return
  if (editingId.value) {
    await update('depart', { id: editingId.value }, { name: form.value.name })
  } else {
    await create('depart', form.value)
  }
  showForm.value = false
  editingId.value = null
  loadData()
}

async function del(d: any) {
  if (!confirm(`确定删除部门 ${d.name}？`)) return
  await remove('depart', { id: d.id })
  loadData()
}

onMounted(() => loadData())
</script>
