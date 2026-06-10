import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login } from '../api/auth'
import { loadToken, saveToken, clearAuth, loadUser, saveUser } from '../utils/token'
import type { UserFormatter } from '../api/types'

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
    // 1. login() sends JSON POST /api/emllogin, returns HttpResponse<UserFormatter>
    // 2. Axios interceptor saves JWT from Authorization header to localStorage
    const res = await login(email, password)
    // res = { message, code, data: UserFormatter }
    const userData: UserFormatter = res.data
    setAuth(loadToken(), { user_id: userData.id, ...userData })
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
