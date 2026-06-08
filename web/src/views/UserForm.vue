<template>
  <div>
    <h4 class="mb-4">{{ isEdit ? '编辑用户' : '新增用户' }}</h4>
    <div class="card card-table">
      <div class="card-body p-4">
        <form @submit.prevent="save">
          <div class="row g-3">
            <div class="col-md-6">
              <label class="form-label small">邮箱</label>
              <input v-model="form.email" type="email" class="form-control" required>
            </div>
            <div class="col-md-6">
              <label class="form-label small">姓名</label>
              <input v-model="form.name" class="form-control">
            </div>
            <div class="col-md-6">
              <label class="form-label small">{{ isEdit ? '新密码（留空不修改）' : '密码' }}</label>
              <input v-model="form.password" type="password" class="form-control" :required="!isEdit">
            </div>
            <div class="col-md-6">
              <label class="form-label small">手机</label>
              <input v-model="form.mobile" class="form-control">
            </div>
            <div class="col-md-4">
              <label class="form-label small">性别</label>
              <select v-model="form.sex" class="form-select">
                <option :value="2">未知</option><option :value="1">男</option><option :value="0">女</option>
              </select>
            </div>
            <div class="col-md-4">
              <label class="form-label small">状态</label>
              <select v-model="form.state" class="form-select">
                <option :value="1">启用</option><option :value="0">禁用</option>
              </select>
            </div>
            <div class="col-md-4">
              <label class="form-label small">别名</label>
              <input v-model="form.aliasname" class="form-control">
            </div>
          </div>
          <div v-if="error" class="alert alert-danger mt-3 py-2 small">{{ error }}</div>
          <div class="mt-4">
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? '保存中...' : '保存' }}
            </button>
            <router-link to="/users" class="btn btn-outline-secondary ms-2">取消</router-link>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { query, create, update } from '../api/mif'

const route = useRoute()
const router = useRouter()
const isEdit = !!route.params.id
const saving = ref(false)
const error = ref('')

const form = ref<any>({
  email: '', name: '', password: '', mobile: '',
  sex: 2, state: 1, aliasname: '', id: '',
})

onMounted(async () => {
  if (isEdit) {
    const res = await query({ Target: 'user', Where: { id: route.params.id }, PageSize: 1 })
    const u = res.data.Data?.[0]
    if (u) {
      form.value = { ...u, password: '' }
    }
  }
})

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (isEdit) {
      const fields: Record<string, any> = { name: form.value.name, mobile: form.value.mobile, sex: form.value.sex, state: form.value.state, aliasname: form.value.aliasname }
      if (form.value.password) fields.password = form.value.password
      await update('user', { id: form.value.id }, fields)
    } else {
      await create('user', form.value)
    }
    router.push('/users')
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>
