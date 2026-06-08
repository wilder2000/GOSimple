import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<any>(JSON.parse(localStorage.getItem('user') || '{}'))

  function setAuth(t: string, u: any) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
    localStorage.setItem('user', JSON.stringify(u))
  }

  async function doLogin(email: string, password: string) {
    const res = await login(email, password)
    setAuth(res.data.token, { user_id: res.data.user_id })
    return res
  }

  function logout() {
    token.value = ''
    user.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function isLoggedIn() {
    return !!token.value
  }

  return { token, user, setAuth, doLogin, logout, isLoggedIn }
})
