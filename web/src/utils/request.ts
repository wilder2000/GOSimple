import axios from 'axios'
import { loadToken, clearAuth } from './token'
import { ECode } from '../api/codes'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

request.interceptors.request.use((config) => {
  const token = loadToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (res) => {
    if (res.status === 401) {
      clearAuth()
      window.location.href = '/admin/login'
      return res
    }
    if (res.status === 200 && res.headers.authorization) {
      localStorage.setItem('token', res.headers.authorization)
    }
    if (res.data?.code === ECode.RNeedLogin) {
      clearAuth()
      window.location.href = '/admin/login'
      return Promise.reject(new Error(res.data.message || '未登录'))
    }
    return res
  },
  (err) => {
    if (err.response?.status === 401) {
      clearAuth()
      window.location.href = '/admin/login'
    }
    return Promise.reject(err)
  },
)

export async function Get<T = any>(url: string, params?: any): Promise<T> {
  const res = await request.get(url, { params })
  return res.data
}

export async function PostJson<T = any>(url: string, data: any): Promise<T> {
  const res = await request.post(url, data, {
    headers: { 'content-type': 'application/json' },
  })
  if (res.data.code !== ECode.SUCCESS) {
    throw new Error(res.data.message || '请求失败')
  }
  return res.data
}

export async function PostForm<T = any>(url: string, data: any): Promise<T> {
  const res = await request.post(url, data, {
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
  })
  if (res.data.code !== ECode.SUCCESS) {
    throw new Error(res.data.message || '请求失败')
  }
  return res.data
}

export function PostFormData(url: string, data: FormData) {
  const config = {
    headers: {
      'Content-Type': 'multipart/form-data',
      Authorization: `Bearer ${loadToken()}`,
    },
  }
  const ax = axios.create()
  return ax.post(url, data, config).then((res) => res.data)
}

export default request
