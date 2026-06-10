import { PostJson, Get } from '../utils/request'
import type { QResponse } from './types'
import type { QRequest, ARequest, URequest, DRequest } from './types'
import { ECmd } from './codes'

// ── Generic CRUD wrappers (MIF pattern) ─────────────────────

type QueryInput = QRequest

export function query<T = any>(req: QueryInput) {
  return PostJson<QResponse<T>>('/v1/mif/q', {
    Target: req.Target,
    PageIndex: req.PageIndex ?? 1,
    PageSize: req.PageSize ?? 15,
    Where: req.Where,
    Order: req.Order,
    Fields: req.Fields,
    Code: req.Code ?? 0,
  })
}

export function create(target: string, data: Record<string, any>) {
  const req: ARequest = { Target: target, ObjectString: JSON.stringify(data) }
  return PostJson('/v1/mif/c', req)
}

export function update(target: string, where: Record<string, any>, fields: Record<string, any>) {
  const req: URequest = { Target: target, Where: where, Fields: fields }
  return PostJson('/v1/mif/u', req)
}

export function remove(target: string, where: Record<string, any>) {
  const req: DRequest = { Target: target, Where: where }
  return PostJson('/v1/mif/d', req)
}

export function queryAll(target: string): Promise<any[]> {
  return query({ Target: target, PageSize: 9999 }).then((res) => res.data.Data || [])
}

// ── MIF "union" queries (multi-target, with selected/unselected) ──

function unionQuery<T = any>(cmd: ECmd, target: string, where: Record<string, any>) {
  const req: QRequest = { Code: cmd, Target: target, Where: where, PageIndex: 1, PageSize: 999 }
  return PostJson<QResponse<T>>('/v1/mif/q', req)
}

// ── User groups (Cmd 1200) ──────────────────────────────────
export function queryUserGroups<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.UserGroups, 'group', where)
}

// ── Groups (Cmd 1201) ───────────────────────────────────────
export function queryGroups<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.QueryGroups, 'group', where)
}

// ── Union groups (Cmd 1202) ─────────────────────────────────
export function queryUnionGroups<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.QueryUnionGroups, 'usergroup', where)
}

// ── Union users (Cmd 1203) ──────────────────────────────────
export function queryUnionUsers<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.QueryUnionUsers, 'usergroup', where)
}

// ── Union roles (Cmd 1204) ──────────────────────────────────
export function queryUnionRoles<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.QueryUnionRoles, 'rolegroup', where)
}

// ── Operators (Cmd 1205) ────────────────────────────────────
export function queryOperators<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.QueryUnionOperators, 'roleoper', where)
}

// ── Departments (Cmd 1206) ──────────────────────────────────
export function queryUnionDepartments<T = any>(where: Record<string, any>) {
  return unionQuery<T>(ECmd.QueryUnionDepartments, 'depuser', where)
}

export { ECmd }
