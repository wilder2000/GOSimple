export class QueryRequest {
  PageIndex = 1
  PageSize = 15
  TotalPages = 0
  Code = 0
  Target = ''
  Order = ''
  Fields: string[] = []
  Where: Record<string, any> = {}

  constructor(cmd: number, target: string) {
    this.Code = cmd
    this.Target = target
  }

  AddPara(name: string, value: any) {
    this.Where[name] = value
  }

  AddLike(name: string, value: string) {
    if (value) {
      this.Where[`${name} like`] = `%${value}%`
    }
  }

  Clear() {
    this.Where = {}
    this.PageIndex = 1
    this.PageSize = 15
    this.TotalPages = 0
    this.Order = ''
    this.Fields = []
  }
}

export class UpdateRequest {
  Code = 0
  Target = ''
  Where: Record<string, any> = {}
  Fields: Record<string, any> = {}

  constructor(cmd: number, target: string) {
    this.Code = cmd
    this.Target = target
  }

  AddWhere(name: string, value: any) {
    this.Where[name] = value
  }

  AddField(name: string, value: any) {
    this.Fields[name] = value
  }

  Clear() {
    this.Where = {}
    this.Fields = {}
  }
}

export class DelRequest {
  Target = ''
  Code = 0
  Where: Record<string, any> = {}

  constructor(cmd: number, target: string) {
    this.Code = cmd
    this.Target = target
  }

  AddPara(name: string, value: any) {
    this.Where[name] = value
  }

  Clear() {
    this.Where = {}
  }
}

export class JSONCreateRequest {
  Target = ''
  ObjectString = ''

  constructor(target: string, data: any) {
    this.Target = target
    this.ObjectString = JSON.stringify(data)
  }
}
