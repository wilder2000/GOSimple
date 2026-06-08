<template>
  <div>
    <h4 class="mb-4">{{ isEdit ? '编辑角色' : '新增角色' }}</h4>
    <div class="card card-table">
      <div class="card-body p-4">
        <form @submit.prevent="save">
          <div class="mb-3">
            <label class="form-label small">角色名称</label>
            <input v-model="form.name" class="form-control" required>
          </div>
          <div v-if="isEdit" class="mb-3">
            <label class="form-label small">关联操作权限</label>
            <div v-if="opers.length" class="border rounded p-3" style="max-height:300px;overflow-y:auto">
              <div v-for="o in opers" :key="o.id" class="form-check form-check-inline">
                <input :id="'op-'+o.id" type="checkbox" class="form-check-input" :value="o.id" v-model="selectedOpers">
                <label :for="'op-'+o.id" class="form-check-label small">{{ o.name }}</label>
              </div>
            </div>
            <small v-else class="text-muted">暂无操作权限数据</small>
          </div>
          <div v-if="error" class="alert alert-danger py-2 small">{{ error }}</div>
          <div class="mt-4">
            <button type="submit" class="btn btn-primary" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
            <router-link to="/roles" class="btn btn-outline-secondary ms-2">取消</router-link>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { query, create, update, remove } from '../api/mif'
import request from '../utils/request'

const route = useRoute()
const router = useRouter()
const isEdit = !!route.params.id
const saving = ref(false)
const error = ref('')
const form = ref<any>({ name: '', id: 0 })
const opers = ref<any[]>([])
const selectedOpers = ref<number[]>([])

onMounted(async () => {
  if (isEdit) {
    const res = await query({ Target: 'role', Where: { id: Number(route.params.id) }, PageSize: 1 })
    const r = res.data.Data?.[0]
    if (r) form.value = { ...r }
    await loadOperators()
  }
})

async function loadOperators() {
  const r = await request.post<any, any>('/v1/rm', {
    Code: 1205, Target: 'roleoper',
    Where: { rid: Number(route.params.id), name: '' },
    PageIndex: 1, PageSize: 999,
  })
  opers.value = r.data?.Data || []
  selectedOpers.value = opers.value.filter((o: any) => o.selected).map((o: any) => o.id)
}

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (isEdit) {
      await update('role', { id: form.value.id }, { name: form.value.name })
      // sync operator-role relations
      const existing = opers.value.filter((o: any) => o.selected).map((o: any) => o.id)
      const toRemove = existing.filter((id: number) => !selectedOpers.value.includes(id))
      const toAdd = selectedOpers.value.filter((id: number) => !existing.includes(id))
      for (const oid of toRemove) {
        await remove('roleoper', { roleid: form.value.id, operatorid: oid })
      }
      for (const oid of toAdd) {
        await create('roleoper', { roleid: form.value.id, operatorid: oid, acte: false })
      }
    } else {
      await create('role', form.value)
    }
    router.push('/roles')
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>
