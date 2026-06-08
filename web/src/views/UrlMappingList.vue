<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="mb-0">URL 映射管理</h4>
      <button class="btn btn-primary btn-sm" @click="showForm = true; form = { operatorid: '', url: '' }"><i class="bi bi-plus-lg"></i> 新增</button>
    </div>
    <div class="card card-table">
      <div class="card-header d-flex justify-content-between">
        <small class="text-muted">共 {{ totalRows }} 条</small>
      </div>
      <div class="table-responsive">
        <table class="table table-hover mb-0">
          <thead class="table-light"><tr><th>ID</th><th>OperatorID</th><th>URL</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="m in list" :key="m.id">
              <td>{{ m.id }}</td><td>{{ m.operatorid }}</td><td><code>{{ m.url }}</code></td>
              <td>
                <button class="btn btn-sm btn-outline-primary me-1" @click="edit(m)"><i class="bi bi-pencil"></i></button>
                <button class="btn btn-sm btn-outline-danger" @click="del(m)"><i class="bi bi-trash"></i></button>
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
          <div class="modal-header"><h6 class="mb-0">{{ editingId ? '编辑映射' : '新增映射' }}</h6><button type="button" class="btn-close" @click="showForm = false"></button></div>
          <div class="modal-body">
            <div class="mb-2">
              <label class="form-label small">OperatorID</label>
              <select v-model.number="form.operatorid" class="form-select">
                <option v-for="o in opers" :key="o.id" :value="o.id">{{ o.id }} - {{ o.name }}</option>
              </select>
            </div>
            <div>
              <label class="form-label small">URL 路径</label>
              <input v-model="form.url" class="form-control" placeholder="/mif/q" @keyup.enter="saveUrl">
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-sm btn-secondary" @click="showForm = false">取消</button>
            <button class="btn btn-sm btn-primary" @click="saveUrl">保存</button>
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
const opers = ref<any[]>([])
const page = ref(1)
const totalPages = ref(0)
const totalRows = ref(0)
const showForm = ref(false)
const form = ref<any>({ operatorid: '', url: '' })
const editingId = ref<number | null>(null)

async function loadData(p?: number) {
  if (p) page.value = p
  const res = await query({ Target: 'urlmap', PageIndex: page.value, PageSize: 15 })
  list.value = res.data.Data || []
  totalPages.value = res.data.TotalPages
  totalRows.value = res.data.TotalRows
}

async function loadOpers() {
  const res = await query({ Target: 'operator', PageSize: 9999 })
  opers.value = res.data.Data || []
}

function edit(m: any) {
  editingId.value = m.id
  form.value = { operatorid: m.operatorid, url: m.url }
  showForm.value = true
}

async function saveUrl() {
  if (!form.value.operatorid || !form.value.url) return
  if (editingId.value) {
    await update('urlmap', { id: editingId.value }, form.value)
  } else {
    await create('urlmap', form.value)
  }
  showForm.value = false
  editingId.value = null
  loadData()
}

async function del(m: any) {
  if (!confirm('确定删除此映射？')) return
  await remove('urlmap', { id: m.id })
  loadData()
}

onMounted(() => { loadData(); loadOpers() })
</script>
