// ============================================================
// Go → TS type mappings for GOSimple HTTP API
// Mirrors Go structs in http/ package
// ============================================================

import type { HttpResponse, DataResponse } from '../types/response'

// ── Common response wrappers ────────────────────────────────
export type { HttpResponse, DataResponse }
export type QResponse<T = any> = HttpResponse<DataResponse<T>>
export type SuccessResponse<T = any> = HttpResponse<T>

// ── Auth / Login ────────────────────────────────────────────
export interface LoginInput {
  email: string
  password: string
}

// Login response: body contains UserFormatter, token is in Authorization header
export type LoginResponse = HttpResponse<UserFormatter>
export type LoginResult = UserFormatter

export interface RegisterUserInput {
  email: string
  password: string
}

export interface RegistUserInput {
  accesskey: string
  secretkey: string
  uuid: string
  email?: string
  name?: string
}

export interface LoginExistInput {
  accesskey?: string
  secretkey?: string
  uuid: string
}

export interface CheckEmailInput {
  email: string
}

// ── User management ─────────────────────────────────────────
export interface RequestUserInput {
  uuid: string
}

export interface DeleteAccountInput {
  uuid: string
}

export interface UpdateAliasInput {
  uuid: string
  aliasname: string
}

export interface ErrorsInput {
  uuid: string
  envinfo: string
  detail: string
}

export interface UserAvatarForm {
  title: string
  email: string
}

// ── Password ────────────────────────────────────────────────
export interface ChangePWD {
  Email: string
  Password: string
  RePassword: string
}

export interface CheckPWD {
  Email: string
  Password: string
}

// ── MIF CRUD ────────────────────────────────────────────────
export interface QRequest {
  PageIndex?: number
  PageSize?: number
  Code?: number
  Target?: string
  Order?: string
  Fields?: string[]
  Attach?: boolean
  Where?: Record<string, any>
}

export interface ARequest {
  Target: string
  Fields?: Record<string, string>
  ObjectString?: string
}

export interface URequest {
  Target: string
  Fields?: Record<string, any>
  Where?: Record<string, any>
}

export interface DRequest {
  Target: string
  Where?: Record<string, any>
}

// ── Role & Permission ───────────────────────────────────────
export interface AddRoleRequest {
  name: string
}

export interface GetRequest {
  Code: number
  FilterName: string
  FilterValue: string
}

export interface UpdateRequest {
  Code: number
  Fields?: Record<string, any>
  Where?: Record<string, any>
}

export interface DeleteRequest {
  Code: number
  Where?: Record<string, any>
}

// ── User formatter (returned by uquery) ─────────────────────
export interface UserFormatter {
  id: string
  name: string
  department: any[]
  email: string
  icon: string
  sex: number
  aliasname: string
  mobile: string
  password: string
  state: number
  reporterror: boolean
}
