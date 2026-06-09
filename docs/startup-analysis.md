# GOSimple 启动过程全链路分析

## 阶段一：Go 包级 `init()` 链（按依赖顺序自动执行）

```
┌─────────────────────────────────────────────────────────────────────┐
│ 1. config.init()                    [config/config.go:101]         │
│    └─ LoadConfig() → Viper 读取 conf/Application.yaml              │
│        → 填充全局 config.AConfig                                     │
│    └─ 同时加载 SysInfo.go (Logo ASCII 常量、环境变量辅助函数)          │
├─────────────────────────────────────────────────────────────────────┤
│ 2. glog.init()                     [glog/GoLogger.go:23]           │
│    └─ 读取 conf/log4g.yaml → 初始化 zap.Logger + lumberjack 轮转    │
│    └─ 设置全局 glog.Logger、glog.LConfig                             │
├─────────────────────────────────────────────────────────────────────┤
│ 3. database.init()                 [database/database-manager.go:15]│
│    └─ LoadDatabaseConfig()                                          │
│    └─ 根据 DataSource.type 选择:                                    │
│       ├─ mysql → gorm.Open(mysql.New, DSN)                         │
│       └─ sqlite → gorm.Open(sqlite.Open, DSN) + WAL 模式            │
│    └─ 设置全局 database.DBHander                                    │
├─────────────────────────────────────────────────────────────────────┤
│ 4. http/validator-config.init()   [http/validator-config.go:23]    │
│    └─ initTrans("zh")                                              │
│    └─ 注册 gin 验证器的中文翻译器（字段名从 json tag 读取）            │
├─────────────────────────────────────────────────────────────────────┤
│ 5. http/url-sync.init()           [http/url-sync.go:21]            │
│    └─ 初始化 urlEntries map (注册表)                                 │
├─────────────────────────────────────────────────────────────────────┤
│ 6. http/mif-initial.init()        [http/mif-initial.go:20]         │
│    └─ RegObject[dbmodel.*]("name") 注册 11 个通用 CRUD 目标:        │
│       user, role, group, groupuser, rolegroup, operator,           │
│       roleoper, depart, depuser, log, urlmap                       │
│    └─ 为 "user" 注入 TargetPreCreateFunc (密码 bcrypt 加密)          │
│    └─ RegisterURL(REQ_{C,Q,U,D}, OPER_ID_ADMIN)                    │
│    └─ AttachMgr.InitHome() → 创建 attach 临时目录                   │
├─────────────────────────────────────────────────────────────────────┤
│ 7. http/http-server.init()        [http/http-server.go:21]         │
│    └─ 打印 Logo + LogoTitle ASCII 横幅                              │
│    └─ 打印 "GOSimple Admin: http://localhost:8090/admin"           │
│    └─ RegMapping[T]() 注册 10 个框架控制器:                          │
│       PwdController      → /pwd    (改密码)                         │
│       CheckPwdController → /cpwd   (验证密码)                       │
│       UserQueryController→ /uquery (用户查询)                       │
│       RoleQueryController→ /rquery (角色查询)                       │
│       RoleAddController  → /radd   (添加角色)                       │
│       UserGroupsController→ /ug    (用户组查询)                     │
│       OperatorController → /rm    (操作权限管理)                    │
│       DepartmentController→ /dm   (部门查询)                        │
│       UserMgrController  → /um    (用户信息修改)                    │
│       UserProfileController→ /upro (用户档案)                       │
└─────────────────────────────────────────────────────────────────────┘
```

## 阶段二：`main()` 执行

```go
// main.go
func main() {
    app.Run(adminHandler())
    //         ↑ 先求值: adminui.go 中的 adminHandler()
    //           读取嵌入的 web/dist → 创建 SPA 文件服务 Handler
}
```

## 阶段三：`app.Run(handler)` 执行

```go
// app/app.go
func Run(adminHandler ...gin.HandlerFunc) {
    // 1. 打印 "GOSimple"
    // 2. 解析 CLI 参数:
    //
    //    ┌─ "-install YES" ─────────────────────────────────┐
    //    │  http.CreateHttpServer(DataSource.Type)          │
    //    │  hs.Install() → startWebServer(install=true)    │
    //    │  → dbscript.Install(UrlMappingsGrouped())       │
    //    │    ├─ 解析嵌入的 initdb.sql → 执行 CREATE TABLE  │
    //    │    ├─ 插入 URL 映射记录                           │
    //    │    └─ 插入种子数据（管理员、角色、组等）            │
    //    │  return                                          │
    //    └─────────────────────────────────────────────────┘
    //
    //    ┌─ "-sync-urls" ───────────────────────────────────┐
    //    │  http.SyncUrlMappings()                          │
    //    │  → 扫描所有已注册的 URL，缺失的插入 s_urlmappings  │
    //    │  return                                          │
    //    └─────────────────────────────────────────────────┘
    //
    // 3. 正常启动:
    hs := http.CreateHttpServer(config.AConfig.Port)   // ":8090"
    if len(adminHandler) > 0 {
        hs.AdminHandler = adminHandler[0]              // 管理后台 Handler
    }
    hs.Start()
}
```

## 阶段四：`hs.Start()` → `startWebServer()`

```go
// http/http-service.go:30
func (hs *HttpServer) Start() {
    startWebServer(hs.Address,        // ":8090"
        *hs.Config.ReadTimeout,       // 30s
        *hs.Config.WriteTimeout,      // 20s
        hs.Actions,                   // 扩展控制器（空）
        hs.NoAuthActions,             // 免认证控制器（空）
        false,                        // install=false
        hs.AdminHandler)              // SPA 文件服务
}
```

## 阶段五：`startWebServer()` 完整流程

```
startWebServer()
│
├─ 1. gin.Default() → 创建 Router（自带 Logger + Recovery 中间件）
│
├─ 2. 挂载静态目录:
│   ├─ Web 映射（如 /out/vip → 对应本地目录）
│   └─ StaticDir（/s → AbsoluteFileDir）
│
├─ 3. 挂载管理后台 SPA（如有）:
│   ├─ router.Any("/admin", adminHandler)
│   └─ router.Any("/admin/*filepath", adminHandler)
│
├─ 4. 注册 Gin 中间件:
│   ├─ 自定义日志格式
│   ├─ Recovery
│   └─ [可选] CORS（AccessControlAllowOrigin=true 时）
│
├─ 5. initController(router) — 注册所有路由:
│   │
│   ├─ 无认证路由（/api/...）:
│   │   ├─ POST /api/emllogin    邮箱密码登录
│   │   ├─ POST /api/reguser     用户注册
│   │   ├─ POST /api/reqmcode    请求验证码
│   │   ├─ POST /api/token_valid 令牌验证
│   │   ├─ POST /api/updmobile   更新手机号
│   │   ├─ POST /api/moblogin    手机号登录
│   │   ├─ POST /api/newreglogin UID注册登录
│   │   └─ POST /api/loginexist  UID登录
│   │
│   ├─ 认证路由 /api/v1（PreProcess 中间件）:
│   │   │
│   │   │  PreProcess 中间件逻辑:
│   │   │  ├─ ParseHttpRequest() → 解析 JWT（Bearer token）
│   │   │  ├─ 验证 token → 提取 userID
│   │   │  ├─ UserAllUrlList(userID) → 查 RBAC 权限链
│   │   │  │    user → groupuser → role → roleoperator → urlmappings
│   │   │  ├─ 检查请求 URL 是否在权限列表中
│   │   │  ├─ 失败计数器（MaxTryTimes = 10, ForbidAccessTime = 5min）
│   │   │  └─ 通过 → c.Next() / 拒绝 → 403 Abort
│   │   │
│   │   ├─ [RegMapping 注册的控制器] → POST /api/v1/{path}
│   │   │   /pwd, /cpwd, /uquery, /rquery, /radd,
│   │   │   /ug, /rm, /dm, /um, /upro
│   │   │
│   │   ├─ [用户管理] → POST /api/v1/{path}
│   │   │   /avatorup, /requestuser, /delaccount,
│   │   │   /modalias, /reperror
│   │   │
│   │   └─ [通用 CRUD] → POST /api/v1/mif/{c,q,u,d}
│   │       └─ 根据 Target 字段分发到 RegObject 注册的 handler
│   │
│   └─ [HttpServer.Actions 扩展] → 挂载到 /api/v1 组
│
├─ 6. SyncUrlMappings() — 扫描所有已注册 URL → 写入 DB
│
├─ 7. 创建 http.Server:
│   ├─ Addr: ":8090"
│   ├─ Handler: gin.Engine
│   ├─ ReadTimeout: 30s
│   └─ WriteTimeout: 20s
│
├─ 8. go srv.ListenAndServe() — 异步启动 HTTP 监听
│
└─ 9. 阻塞等待信号:
    ├─ signal.Notify(SIGINT, SIGTERM)
    └─ 收到信号 → srv.Shutdown(5s timeout) → 优雅退出
```

## 启动时序总图

```
时间 →
─────┬─────────────────────────────────────────────────────────
     │ Go 程序启动
     │
     ├─ config.init()        ─── 加载配置 → AConfig
     ├─ glog.init()          ─── 初始化日志 → Logger
     ├─ database.init()      ─── 连接数据库 → DBHander
     ├─ validator.init()     ─── 注册中文验证器
     ├─ url-sync.init()      ─── 初始化 URL 注册表
     ├─ mif-initial.init()   ─── 注册 11 个 CRUD 目标
     ├─ http-server.init()   ─── 注册 10 个框架控制器 + 打印 Logo
     │
     ├─ main()
     │   ├─ adminHandler()   ─── 读取嵌入的 web/dist → 创建 SPA handler
     │   └─ app.Run()
     │       ├─ 解析 CLI 参数
     │       ├─ CreateHttpServer(:8090)
     │       ├─ 设置 AdminHandler
     │       └─ Start()
     │           └─ startWebServer()
     │               ├─ gin.Default()
     │               ├─ 挂载静态文件
     │               ├─ 挂载 Admin SPA
     │               ├─ 注册 Middleware（CORS）
     │               ├─ initController()
     │               │   ├─ 无认证路由
     │               │   ├─ /api/v1 组 + PreProcess
     │               │   ├─ 框架控制器
     │               │   ├─ 用户控制器
     │               │   └─ MIF 通用 CRUD
     │               ├─ SyncUrlMappings()
     │               ├─ srv.ListenAndServe()  ← 服务就绪
     │               └─ 等待信号 ← 阻塞在这里
```

## 关键设计点

| 特性 | 说明 |
|------|------|
| **全局单例** | `config.AConfig`、`database.DBHander`、`glog.Logger` 均由 `init()` 在 `main()` 之前填充 |
| **自注册机制** | 控制器通过 `RegMapping[T]()` 在 `init()` 中注册到全局 `mappings` map |
| **泛型 CRUD** | `RegObject[T]("name")` 用反射注册通用增删改查，`mif-controller.go` 中统一分发 |
| **RBAC** | `PreProcess` 中间件在每个认证请求时做 JWT 解析 + URL 权限校验 |
| **URL 同步** | 每次启动自动将代码中注册的 URL 写入 `s_urlmappings` 表，无需手动维护 |
| **优雅关闭** | `signal.Notify` 捕获中断信号 → `Shutdown(5s)` → 完成正在处理的请求后退出 |
