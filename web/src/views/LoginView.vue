<template>
  <div class="login-wrapper">
    <div class="login-card card">
      <div class="card-body">
        <div class="text-center mb-4">
          <div class="brand-text mb-2">GOSimple</div>
          <p class="text-muted small">后台管理系统</p>
        </div>
        <form @submit.prevent="handleLogin">
          <div class="mb-3">
            <label class="form-label small">邮箱</label>
            <input v-model="email" type="email" class="form-control" placeholder="请输入邮箱" required autofocus>
          </div>
          <div class="mb-3">
            <label class="form-label small">密码</label>
            <input v-model="password" type="password" class="form-control" placeholder="请输入密码" required>
          </div>
          <div v-if="error" class="alert alert-danger py-2 small">{{ error }}</div>
          <button type="submit" class="btn btn-primary w-100" :disabled="loading">
            <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>
        <div class="mt-3 text-center">
          <small class="text-muted">默认管理员: wild.shang@163.com / admin@123</small>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.doLogin(email.value, password.value)
    router.push('/dashboard')
  } catch (e: any) {
    error.value = e.message || '登录失败，请检查邮箱和密码'
  } finally {
    loading.value = false
  }
}
</script>
