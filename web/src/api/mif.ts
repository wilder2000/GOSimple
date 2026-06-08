import request from '../utils/request'

export interface PageData<T> {
  PageIndex: number
  PageSize: number
  TotalPages: number
  TotalRows: number
  Data: T[]
}

export interface QueryParams {
  Target: string
  Where?: Record<string, any>
  PageIndex?: number
  PageSize?: number
  Order?: string
  Fields?: string[]
}

export interface QResponse {
  message: string
  code: number
  data: PageData<any>
}

export function query<T = any>(params: QueryParams) {
  return request.post<any, QResponse>('/v1/mif/q', params)
}

export function create(Target: string, data: Record<string, any>) {
  return request.post('/v1/mif/c', {
    Target,
    ObjectString: JSON.stringify(data),
  })
}

export function update(Target: string, where: Record<string, any>, fields: Record<string, any>) {
  return request.post('/v1/mif/u', { Target, Where: where, Fields: fields })
}

export function remove(Target: string, where: Record<string, any>) {
  return request.post('/v1/mif/d', { Target, Where: where })
}

export function queryAll(target: string): Promise<any[]> {
  return query({ Target: target, PageSize: 9999 }).then((res) => res.data.Data || [])
}
