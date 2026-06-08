<template>
  <div>
    <h4 class="mb-4">{{ isEdit ? '编辑编组' : '新增编组' }}</h4>
    <div class="card card-table">
      <div class="card-body p-4">
        <form @submit.prevent="save">
          <div class="mb-3">
            <label class="form-label small">编组名称</label>
            <input v-model="form.name" class="form-control" required>
          </div>
          <div v-if="isEdit" class="row g-3 mb-3">
            <div class="col-md-6">
              <label class="form-label small">关联用户</label>
              <div class="border rounded p-2" style="max-height:250px;overflow-y:auto">
                <div v-for="u in allUsers" :key="u.id" class="form-check">
                  <input :id="'gu-'+u.id" type="checkbox" class="form-check-input" :value="u.id" v-model="selectedUsers">
                  <label :for="'gu-'+u.id" class="form-check-label small">{{ u.email }} ({{ u.name || '-' }})</label>
                </div>
                <div v-if="!allUsers.length" class="text-muted small py-2">暂无用户</div>
              </div>
            </div>
            <div class="col-md-6">
              <label class="form-label small">关联角色</label>
              <div class="border rounded p-2" style="max-height:250px;overflow-y:auto">
                <div v-for="r in allRoles" :key="r.id" class="form-check">
                  <input :id="'gr-'+r.id" type="checkbox" class="form-check-input" :value="r.id" v-model="selectedRoles">
                  <label :for="'gr-'+r.id" class="form-check-label small">{{ r.name }}</label>
                </div>
                <div v-if="!allRoles.length" class="text-muted small py-2">暂无角色</div>
              </div>
            </div>
          </div>
          <div v-if="error" class="alert alert-danger py-2 small">{{ error }}</div>
          <div class="mt-4">
            <button type="submit" class="btn btn-primary" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
            <router-link to="/groups" class="btn btn-outline-secondary ms-2">取消</router-link>
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
const allUsers = ref<any[]>([])
const allRoles = ref<any[]>([])
const selectedUsers = ref<string[]>([])
const selectedRoles = ref<number[]>([])

onMounted(async () => {
  const [users, roles] = await Promise.all([
    query({ Target: 'user', PageSize: 9999 }),
    query({ Target: 'role', PageSize: 9999 }),
  ])
  allUsers.value = users.data.Data || []
  allRoles.value = roles.data.Data || []

  if (isEdit) {
    const res = await query({ Target: 'group', Where: { id: Number(route.params.id) }, PageSize: 1 })
    const g = res.data.Data?.[0]
    if (g) form.value = { ...g }
    // load group-user and group-role relations
    const [gu, gr] = await Promise.all([
      query({ Target: 'groupuser', Where: { groupid: Number(route.params.id) }, PageSize: 9999 }),
      query({ Target: 'rolegroup', Where: { groupid: Number(route.params.id) }, PageSize: 9999 }),
    ])
    selectedUsers.value = (gu.data.Data || []).map((x: any) => x.userid)
    selectedRoles.value = (gr.data.Data || []).map((x: any) => x.roleid)
  }
})

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (isEdit) {
      await update('group', { id: form.value.id }, { name: form.value.name })
      // sync users
      const existingU = (await query({ Target: 'groupuser', Where: { groupid: form.value.id }, PageSize: 9999 })).data.Data || []
      for (const x of existingU) { if (!selectedUsers.value.includes(x.userid)) await remove('groupuser', { groupid: form.value.id, userid: x.userid }) }
      for (const uid of selectedUsers.value) { if (!existingU.find((x: any) => x.userid === uid)) await create('groupuser', { groupid: form.value.id, userid: uid }) }
      // sync roles
      const existingR = (await query({ Target: 'rolegroup', Where: { groupid: form.value.id }, PageSize: 9999 })).data.Data || []
      for (const x of existingR) { if (!selectedRoles.value.includes(x.roleid)) await remove('rolegroup', { groupid: form.value.id, roleid: x.roleid }) }
      for (const rid of selectedRoles.value) { if (!existingR.find((x: any) => x.roleid === rid)) await create('rolegroup', { groupid: form.value.id, roleid: rid }) }
    } else {
      await create('group', form.value)
    }
    router.push('/groups')
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>
