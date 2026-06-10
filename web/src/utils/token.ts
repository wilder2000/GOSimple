const TOKEN_KEY = 'token'
const USER_KEY = 'user'

export function loadToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function saveToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

export function loadUser(): any {
  try {
    return JSON.parse(localStorage.getItem(USER_KEY) || '{}')
  } catch {
    return {}
  }
}

export function saveUser(user: any) {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}
