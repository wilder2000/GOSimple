<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="mb-0">功能权限管理</h4>
      <button class="btn btn-primary btn-sm" @click="showForm = true; form = { id: 0, name: '' }"><i class="bi bi-plus-lg"></i> 新增</button>
    </div>
    <div class="card card-table">
      <div class="card-header d-flex justify-content-between"><small class="text-muted">共 {{ totalRows }} 条</small></div>
      <div class="table-responsive">
        <table class="table table-hover mb-0">
          <thead class="table-light"><tr><th>ID</th><th>名称</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="o in list" :key="o.id">
              <td>{{ o.id }}</td><td>{{ o.name }}</td>
              <td>
                <button class="btn btn-sm btn-outline-primary me-1" @click="edit(o)"><i class="bi bi-pencil"></i></button>
                <button class="btn btn-sm btn-outline-danger" @click="del(o)"><i class="bi bi-trash"></i></button>
              </td>
            </tr>
            <tr v-if="!list.length"><td colspan="3" class="text-center text-muted py-4">暂无数据</td></tr>
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
          <div class="modal-header"><h6 class="mb-0">{{ editingId ? '编辑权限' : '新增权限' }}</h6><button type="button" class="btn-close" @click="showForm = false"></button></div>
          <div class="modal-body">
            <div class="mb-2"><label class="form-label small">ID</label><input v-model.number="form.id" type="number" class="form-control" :disabled="!!editingId"></div>
            <div><label class="form-label small">名称</label><input v-model="form.name" class="form-control" placeholder="功能名称" @keyup.enter="saveOp"></div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-sm btn-secondary" @click="showForm = false">取消</button>
            <button class="btn btn-sm btn-primary" @click="saveOp">保存</button>
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
const form = ref<any>({ id: 0, name: '' })
const editingId = ref<number | null>(null)

async function loadData(p?: number) {
  if (p) page.value = p
  const res = await query({ Target: 'operator', PageIndex: page.value, PageSize: 15 })
  list.value = res.data.Data || []
  totalPages.value = res.data.TotalPages
  totalRows.value = res.data.TotalRows
}

function edit(o: any) {
  editingId.value = o.id
  form.value = { ...o }
  showForm.value = true
}

async function saveOp() {
  if (!form.value.name || !form.value.id) return
  if (editingId.value) {
    await update('operator', { id: editingId.value }, { name: form.value.name })
  } else {
    await create('operator', form.value)
  }
  showForm.value = false
  editingId.value = null
  loadData()
}

async function del(o: any) {
  if (!confirm(`确定删除权限 ${o.name}？`)) return
  await remove('operator', { id: o.id })
  loadData()
}

onMounted(() => loadData())
</script>
