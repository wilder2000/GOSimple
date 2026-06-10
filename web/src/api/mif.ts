import { PostJson } from '../utils/request'
import type { HttpResponse, DataResponse } from '../types/response'
import { QueryRequest, UpdateRequest, DelRequest, JSONCreateRequest } from './requests'

export type QResponse<T = any> = HttpResponse<DataResponse<T>>

type QueryInput = QueryRequest | {
  Target: string
  Where?: Record<string, any>
  PageIndex?: number
  PageSize?: number
  Order?: string
  Fields?: string[]
}

export function query<T = any>(req: QueryInput) {
  const payload = req instanceof QueryRequest ? req : {
    Target: req.Target,
    PageIndex: req.PageIndex ?? 1,
    PageSize: req.PageSize ?? 15,
    Where: req.Where,
    Order: req.Order,
    Fields: req.Fields,
  }
  return PostJson<QResponse<T>>('/v1/mif/q', payload)
}

export function create(target: string, data: Record<string, any>) {
  const req = new JSONCreateRequest(target, data)
  return PostJson('/v1/mif/c', req)
}

export function update(target: string, where: Record<string, any>, fields: Record<string, any>) {
  const req = new UpdateRequest(0, target)
  req.Where = where
  req.Fields = fields
  return PostJson('/v1/mif/u', req)
}

export function remove(target: string, where: Record<string, any>) {
  const req = new DelRequest(0, target)
  req.Where = where
  return PostJson('/v1/mif/d', req)
}

export function queryAll(target: string): Promise<any[]> {
  const req = new QueryRequest(0, target)
  req.PageSize = 9999
  return PostJson<QResponse>('/v1/mif/q', req).then((res) => res.data.Data || [])
}
