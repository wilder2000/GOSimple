// ── Response codes ──────────────────────────────────────────
export enum ECode {
  SUCCESS = 0,
  RFailed = 1,
  RNeedLogin = 2,
  RServerError = 3,
  LoginParaFormat = 21,
  PwdWrong = 22,
  UserNotFound = 23,
  DataExistFound = 24,
  ParaFormat = 25,
}

// ── Command codes (mirrors Go Cmd* constants) ───────────────
export enum ECmd {
  ImagePrefixPage = 1100,

  UserGroups            = 1200,
  QueryGroups           = 1201,
  QueryUnionGroups      = 1202,
  QueryUnionUsers       = 1203,
  QueryUnionRoles       = 1204,
  QueryUnionOperators   = 1205,
  QueryUnionDepartments = 1206,

  UpdateUserName = 1207,
  UpdateUserSex  = 1208,
  GetUserBasic   = 1209,
}

// ── User state ──────────────────────────────────────────────
export enum EUserState {
  Disable = 0,
  Enable  = 1,
}
