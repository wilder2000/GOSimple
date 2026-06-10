import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login } from '../api/auth'
import { loadToken, saveToken, clearAuth, loadUser, saveUser } from '../utils/token'
import type { LoginResult } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(loadToken())
  const user = ref<any>(loadUser())

  function setAuth(t: string, u: any) {
    token.value = t
    user.value = u
    saveToken(t)
    saveUser(u)
  }

  async function doLogin(email: string, password: string) {
    const res = await login(email, password)
    setAuth(res.data.token, { user_id: res.data.user_id })
    return res
  }

  function logout() {
    token.value = ''
    user.value = {}
    clearAuth()
  }

  function isLoggedIn() {
    return !!token.value
  }

  return { token, user, setAuth, doLogin, logout, isLoggedIn }
})
