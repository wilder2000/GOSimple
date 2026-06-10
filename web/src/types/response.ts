export interface HttpResponse<T> {
  message: string
  code: number
  data: T
}

export interface DataResponse<T> {
  PageIndex: number
  PageSize: number
  TotalPages: number
  TotalRows: number
  Data: T[]
}
