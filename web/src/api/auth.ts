import request from '../utils/request'
import { saveToken } from '../utils/token'

export interface LoginResult {
  message: string
  code: number
  data: {
    user_id: string
    token: string
  }
}

export function login(email: string, password: string) {
  const params = new URLSearchParams()
  params.append('Email', email)
  params.append('Password', password)
  return request.post<any, LoginResult>('/emllogin', params, {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  })
}

export function handleLoginResult(res: LoginResult) {
  if (res.data?.token) {
    saveToken(res.data.token)
  }
}
