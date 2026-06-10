import request, { PostJson, PostForm, PostFormData } from '../utils/request'
import { saveToken } from '../utils/token'
import type { HttpResponse, RegisterUserInput, RegistUserInput, ChangePWD, CheckPWD, ErrorsInput, SuccessResponse, LoginResult } from './types'
export type { LoginResult }

// ── Login (form-urlencoded) ─────────────────────────────────
export function login(email: string, password: string) {
  const params = new URLSearchParams()
  params.append('Email', email)
  params.append('Password', password)
  return request.post<any, HttpResponse<LoginResult>>('/emllogin', params, {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  })
}

export function handleLoginResult(res: HttpResponse<LoginResult>) {
  if (res.data?.token) {
    saveToken(res.data.token)
  }
}

// ── Register (JSON) ─────────────────────────────────────────
export function registerUser(data: RegisterUserInput) {
  return PostJson<SuccessResponse>('/reguser', data)
}

// ── Request mobile code (form) ──────────────────────────────
export function requestMobileCode(mobile: string) {
  const params = new URLSearchParams()
  params.append('mobile', mobile)
  return PostForm<SuccessResponse>('/reqmcode', params)
}

// ── Update mobile (form) ────────────────────────────────────
export function updateMobile(mobile: string, vcode: string, uid: string) {
  const params = new URLSearchParams()
  params.append('mobile', mobile)
  params.append('vcode', vcode)
  params.append('uid', uid)
  return PostForm<SuccessResponse>('/updmobile', params)
}

// ── Mobile login (form) ─────────────────────────────────────
export function mobileLogin(mobile: string, vcode: string) {
  const params = new URLSearchParams()
  params.append('mobile', mobile)
  params.append('vcode', vcode)
  return PostForm<HttpResponse<LoginResult>>('/moblogin', params)
}

// ── New registration with UID (JSON) ────────────────────────
export function uidLoginRegist(data: RegistUserInput) {
  return PostJson<HttpResponse<LoginResult>>('/newreglogin', data)
}

// ── Login with existing UID (JSON) ──────────────────────────
export function uidLoginWithExist(data: { uuid: string; accesskey?: string; secretkey?: string }) {
  return PostJson<HttpResponse<LoginResult>>('/loginexist', data)
}

// ── Change password (auth) ──────────────────────────────────
export function changePassword(data: ChangePWD) {
  return PostJson<SuccessResponse>('/v1/pwd', data)
}

// ── Check password (auth) ───────────────────────────────────
export function checkPassword(data: CheckPWD) {
  return PostJson<SuccessResponse>('/v1/cpwd', data)
}

// ── Upload avatar (auth, multipart) ─────────────────────────
export function uploadAvatar(data: FormData) {
  return PostFormData('/v1/avatorup', data)
}

// ── Request user info (auth) ────────────────────────────────
export function requestUser(uuid: string) {
  return PostJson<HttpResponse<any>>('/v1/requestuser', { uuid })
}

// ── Delete account (auth) ───────────────────────────────────
export function deleteAccount(uuid: string) {
  return PostJson<SuccessResponse>('/v1/delaccount', { uuid })
}

// ── Update alias name (auth) ────────────────────────────────
export function updateAliasName(uuid: string, aliasname: string) {
  return PostJson<SuccessResponse>('/v1/modalias', { uuid, aliasname })
}

// ── Report error (auth) ─────────────────────────────────────
export function reportError(data: ErrorsInput) {
  return PostJson<SuccessResponse>('/v1/reperror', data)
}

// ── Add role (auth) ─────────────────────────────────────────
export function addRole(name: string) {
  return PostJson<SuccessResponse>('/v1/radd', { name })
}
