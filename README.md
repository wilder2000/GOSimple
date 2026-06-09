# GOSimple —— 通用 Go 后端框架

基于 Gin + GORM + JWT 的通用后端框架，内置用户认证、RBAC 权限控制、通用 CRUD、文件导出等功能。

---

## 快速开始

### 环境要求

- Go 1.23.4+
- MySQL / SQLite

### 安装

```bash
# 克隆项目
git clone <repo_url>
cd GOSimple

# 初始化数据库（创建表结构和种子数据）
go run . -install YES
```

> `-install YES` 执行顺序：框架核心表 → URL 映射 → 种子数据 → **模块注册的 SQL**（见[模块安装 SQL](#模块安装-sql)）。
> 模块的自定义表也可以通过 `init()` + `AutoMigrate` 在启动前自动创建（见[自动建表策略](#自动建表策略)）。

### 构建前端（可选）

管理后台使用 Vue 3 构建，首次使用前需要编译前端资源：

```bash
cd web && npm install && npm run build
```

编译后的文件会自动嵌入到 Go 二进制中，无需额外部署。

### 同步 URL 映射（可选）

新增 Controller 后，可手动同步 URL 映射到数据库，无需重新安装：
```bash
go run . -sync-urls
```

> 正常启动服务时也会自动补录缺失的 URL 映射，无需手动干预。

### 启动服务

```bash
go run .
```

服务默认监听 `:8090`，启动后可通过 http://localhost:8090/admin 访问管理后台。

---

## 管理后台

系统内置一套基于 Vue 3 + Bootstrap 5 的权限管理后台，提供完整的 RBAC 管理界面。

### 访问地址

启动服务后访问：**http://localhost:8090/admin**

### 默认管理员账号

| 字段 | 值 |
|------|-----|
| 邮箱 | `wild.shang@163.com` |
| 密码 | `admin@123`（由配置文件 `Security.DefaultAdminPWD` 指定） |

### 功能页面

| 菜单 | 路径 | 说明 |
|------|------|------|
| 仪表盘 | `/admin/dashboard` | 系统概览，展示用户/角色/编组/部门数量 |
| 用户管理 | `/admin/users` | 用户列表、新增、编辑、删除；新增时密码自动 bcrypt 加密 |
| 角色管理 | `/admin/roles` | 角色 CRUD；编辑时可勾选关联的操作权限 |
| 编组管理 | `/admin/groups` | 编组 CRUD；编辑时可勾选关联用户和角色 |
| 部门管理 | `/admin/departments` | 部门 CRUD（弹窗操作） |
| 功能权限 | `/admin/operators` | 操作权限定义 CRUD（弹窗操作） |
| URL 映射 | `/admin/urlmappings` | URL 路径与操作权限的映射关系管理 |
| 登录日志 | `/admin/logs` | 查看用户登录记录（只读） |

### 前端技术栈

| 技术 | 用途 |
|------|------|
| Vue 3 (Composition API + TypeScript) | 前端框架 |
| Vite 6 | 构建工具 |
| Vue Router 4 | SPA 路由（History 模式） |
| Pinia | 状态管理 |
| Axios | HTTP 客户端（JWT 自动注入） |
| Bootstrap 5.3 | CSS 样式布局 |
| Bootstrap Icons | 图标库 |

---

## 配置文件

### conf/Application.yaml

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `Port` | string | `:8090` | HTTP 监听地址 |
| `ReadTimeout` | int | `20` | 读超时（秒） |
| `WriteTimeout` | int | `20` | 写超时（秒） |
| `ExpireTime` | int | `10` | JWT 令牌过期时间（分钟） |
| `DocProcPoolSize` | int | `20` | 文档处理线程池大小 |
| `MaxCurrThread` | int | `3` | 最大并发线程数 |
| `DataSource.dsn` | string | MySQL DSN | 数据库连接字符串 |
| `DataSource.type` | string | `mysql` | 数据库类型：`mysql` / `sqlite` |
| `DataSource.maxidleconnections` | int | `10` | 最大空闲连接数 |
| `DataSource.maxopenconnections` | int | `5` | 最大打开连接数 |
| `Security.Registration` | bool | `true` | 是否允许公开注册 |
| `Security.MaxTryTimes` | int | `10` | 失败尝试次数锁定阈值 |
| `Security.ForbidAccessTime` | float | `5.0` | 锁定时间（分钟） |
| `Security.DefaultAdminPWD` | string | `admin@123` | 安装时默认管理员密码 |
| `AppSecret.AccessKey` | string | `1233333333` | App 访问密钥 |
| `AppSecret.SecretKey` | string | `hsdf11212` | App 密钥 |
| `AccessControlAllowOrigin` | bool | `false` | 是否启用 CORS |
| `StaticDir.RelativePath` | string | `/s` | 静态文件 URL 前缀 |
| `StaticDir.AbsoluteFileDir` | string | `/usr/www` | 静态文件本地路径 |

### 环境变量覆盖

- `GOGO_HOME` — 覆盖应用目录（默认：当前目录）
- `GOGO_CONFIG_FILE` — 覆盖配置文件名称（默认：`Application.yaml`）

### conf/log4g.yaml

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `file` | string | `application.log` | 日志文件名 |
| `dir` | string | `./` | 日志目录 |
| `MaxSize` | int | `10` | 日志文件最大大小（MB） |
| `MaxAge` | int | `30` | 日志保留天数 |
| `MaxBackups` | int | `7` | 保留的旧日志文件数 |
| `Compress` | bool | `true` | 是否压缩归档日志 |
| `LogLevel` | string | `info` | 日志级别：`debug` / `info` / `error` |
| `Console` | bool | `true` | 是否同时输出到控制台 |

---

## API 接口

所有接口统一返回 JSON 格式：`{ "message": "...", "code": 0, "data": ... }`。

### 公开接口（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/emllogin` | 邮箱密码登录（管理后台） |
| POST | `/api/reguser` | 用户注册（需 `Security.Registration=true`） |
| POST | `/api/reqmcode` | 请求手机验证码 |
| POST | `/api/token_valid` | 令牌验证（用于 Nginx auth_request 静态文件鉴权） |
| POST | `/api/updmobile` | 更新手机号（需验证码） |
| POST | `/api/moblogin` | 手机号 + 验证码登录 |
| POST | `/api/newreglogin` | UID 自动注册/登录（App 端） |
| POST | `/api/loginexist` | UID 登录（已注册用户） |

#### 登录接口请求/响应示例

**POST /api/emllogin**

请求：
```json
{ "email": "admin@example.com", "password": "admin@123" }
```

响应：
```json
{
  "message": "success",
  "code": 0,
  "data": {
    "id": "Administrator",
    "name": "流星划过黑暗的夜空つ",
    "email": "wild.shang@163.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

认证令牌通过 `Authorization: Bearer <token>` 响应头返回。

### 认证接口（需要 JWT）

所有认证接口前缀为 `/api/v1`，请求头需携带 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/cpwd` | 验证当前密码 |
| POST | `/api/v1/pwd` | 修改密码 |
| POST | `/api/v1/uquery` | 查询用户列表（分页） |
| POST | `/api/v1/rquery` | 查询角色列表（分页） |
| POST | `/api/v1/radd` | 添加角色 |
| POST | `/api/v1/ug` | 查询用户组及关联关系 |
| POST | `/api/v1/rm` | 角色管理——查询操作权限及关联状态 |
| POST | `/api/v1/dm` | 部门管理——查询部门用户及关联状态 |
| POST | `/api/v1/um` | 修改用户姓名/性别 |
| POST | `/api/v1/upro` | 获取用户基本信息及所属部门 |
| POST | `/api/v1/avatorup` | 上传用户头像 |
| POST | `/api/v1/requestuser` | 根据 UID 查询用户信息 |
| POST | `/api/v1/delaccount` | 删除账号（级联清理） |
| POST | `/api/v1/modalias` | 修改用户别名 |
| POST | `/api/v1/reperror` | 上报客户端错误信息 |

---

## 通用 CRUD 接口（MIF）

系统提供一套泛型通用 CRUD 框架，注册后即可对任意数据模型进行增删改查。

### 注册的目标对象

| 目标名称 (`Target`) | 对应模型 | 数据库表 |
|---------------------|----------|----------|
| `user` | SUser | s_users |
| `role` | SRole | s_role |
| `group` | SGroup | s_group |
| `groupuser` | SGroupuser | s_groupuser |
| `rolegroup` | SRolegroup | s_rolegroup |
| `operator` | SOperator | s_operators |
| `roleoper` | SRoleoperator | s_roleoperator |
| `depart` | SDepartment | s_departments |
| `depuser` | SDepuser | s_depusers |
| `log` | SLog | s_logs |
| `urlmap` | SUrlMapping | s_urlmappings |

### 通用查询

**POST /api/v1/mif/q**

```json
{
  "Target": "user",
  "PageIndex": 1,
  "PageSize": 20,
  "Order": "createtime desc",
  "Fields": ["id", "name", "email"],
  "Where": { "state": 1, "name like": "%张%" },
  "attach": false
}
```

响应：
```json
{
  "message": "success",
  "code": 0,
  "data": {
    "PageIndex": 1,
    "PageSize": 20,
    "TotalPages": 5,
    "Data": [...],
    "attach": ""
  }
}
```

> 设置 `attach: true` 可将查询结果导出为 Excel 文件，返回下载链接。

### 通用创建

**POST /api/v1/mif/c**

```json
{
  "Target": "role",
  "ObjectString": "{\"name\":\"新角色\"}"
}
```

### 通用更新

**POST /api/v1/mif/u**

```json
{
  "Target": "user",
  "Fields": { "name": "新名字" },
  "Where": { "id": "xxx" }
}
```

### 通用删除

**POST /api/v1/mif/d**

```json
{
  "Target": "user",
  "Where": { "id": "xxx" }
}
```

### 查询条件语法

`Where` 中的 key 支持操作符后缀：

| 语法 | 示例 | 说明 |
|------|------|------|
| `field` | `"name": "张三"` | 等值查询 |
| `field =` | `"state =": 1` | 等值查询 |
| `field >` | `"age >": 18` | 大于 |
| `field >=` | `"age >=": 18` | 大于等于 |
| `field <` | `"age <": 60` | 小于 |
| `field <=` | `"age <=": 60` | 小于等于 |
| `field !=` | `"state !=": 0` | 不等于 |
| `field like` | `"name like": "%张%"` | 模糊匹配 |
| `field in` | `"id in": ["a","b"]` | IN 查询 |
| `field` | `"field": "wk_is_null"` | IS NULL |
| `field` | `"field": "wk_is_not_null"` | IS NOT NULL |

---

## 数据库模型

系统包含以下数据库表（位于 `dbscript/MySQL/initdb.sql`）：

| 表名 | 说明 |
|------|------|
| `s_users` | 用户表（ID, 邮箱, 密码, 姓名, 头像, 别名, 手机号, 性别, 状态） |
| `s_role` | 角色表 |
| `s_group` | 用户组表 |
| `s_groupuser` | 用户-组关联表（M:M） |
| `s_rolegroup` | 角色-组关联表（M:M） |
| `s_operators` | 操作权限定义表 |
| `s_roleoperator` | 角色-操作权限关联表（M:M） |
| `s_departments` | 部门表 |
| `s_depusers` | 用户-部门关联表（M:M） |
| `s_resources` | 资源权限表（读写下载） |
| `s_logs` | 登录日志表 |
| `s_items` | 参考数据项表 |
| `s_debug` | 客户端错误报告表 |
| `s_search` | 搜索日志表 |
| `s_urlmappings` | URL 路径-操作权限映射表 |

用户状态（`s_users.state`）：
- `0` — 管理员
- `1` — 普通用户
- `2` — 自动注册
- `3` — VIP
- `4` — SVIP
- `999` — 锁定

---

## 认证与授权

### JWT 认证

- 使用 HS256 签名算法
- 令牌中包含 `user_id`、`exp`（过期时间）、`user_ip` 声明
- 每次请求自动刷新令牌（新令牌通过响应头 `Authorization` 返回）
- 令牌有效期通过配置文件 `ExpireTime` 设置（默认 10 分钟）

### RBAC 权限控制

权限链路：
```
用户 → 用户组 → 角色 → 操作权限 → URL 路径
```

- 每个请求经过 `PreProcess` 中间件，验证用户是否有权访问目标 URL
- 连续 `MaxTryTimes` 次失败访问后，用户会被锁定 `ForbidAccessTime` 分钟

### URL 映射管理

URL 映射将 API 路径与操作权限（operator）关联，决定用户能否访问某个接口。

- **启动自动同步**：服务启动时自动扫描所有已注册的路由，将缺失的 URL 映射插入 `s_urlmappings` 表
- **声明式权限**：Controller 可通过实现 `OperatorProvider` 接口声明自己所属的 operator ID，不声明的默认归管理员（ID=10）
- **运行时同步**：运行 `go run . -sync-urls` 可手动触发 URL 映射同步（不启动 HTTP 服务）

示例——在 Controller 中声明权限归属：

```go
type UserProfileController struct {
    AbstractController[GetRequest]
}

func (s UserProfileController) OperatorId() int32 { return 11 } // 查看者
func (s UserProfileController) UrlPath() string   { return "/upro" }
```

预定义的 operator：

| ID | 名称 | 说明 |
|----|------|------|
| 10 | 管理功能 | 管理员，拥有所有 URL 访问权限 |
| 11 | 查看功能 | 终端用户，仅可访问声明为 `OPER_ID_VIEWER` 的 URL |

### 内置种子数据

安装时自动创建：
- **管理员角色**（id=2）：拥有全部操作权限
- **终端用户角色**（id=1）：仅拥有查看权限
- **管理员组**（id=2）：关联管理员角色
- **终端用户组**（id=1）：关联终端用户角色
- **默认管理员**：ID=`Administrator`，邮箱=`wild.shang@163.com`，密码默认为 `admin@123`

---

## 功能特性

### 1. 日志系统 (glog)

基于 Uber Zap + Lumberjack 实现日志轮转，支持按大小/时间轮转和压缩归档。

```go
import "github.com/wilder2000/GOSimple/glog"

glog.Logger.Info("信息日志")
glog.Logger.InfoF("格式化 %s", "日志")
glog.Logger.Error("错误日志")
glog.Logger.DebugF("调试信息: %v", data)
```

### 2. 线程池 (pool)

泛型线程池，支持超时机制：

```go
import "github.com/wilder2000/GOSimple/pool"

proc := func(task *MyTask) {
    // 处理任务
}
engine := pool.New[MyTask](pooSize, currThread, proc)
engine.Append(&MyTask{...})
```

### 3. Excel 导出

通用查询设置 `attach: true` 即可将查询结果导出为 `.xlsx` 文件，自动生成列标题和下载链接。

### 4. 验证器

使用 `go-playground/validator`，内置中文错误消息翻译。支持 `alias` 标签自定义字段名称。

### 5. 静态文件服务

通过 `StaticDir` 配置映射静态文件目录，可用于存放用户上传的头像、导出文件等。

### 6. 多数据库支持

支持 MySQL 和 SQLite，通过 `DataSource.type` 切换。`config/database.go` 的 `LoadDatabaseConfig()` 可在运行时重新加载数据库配置。

---

## 二次开发

GOSimple 作为独立 Go module（`github.com/wilder2000/GOSimple`）可被外部项目引用。框架提供了可导入的 `app` 包，外部项目只需写一个薄 `main.go` 即可启动。

### 二开项目结构

```
myapp/                              # 你的业务项目（独立 Go module）
├── go.mod                          # require github.com/wilder2000/GOSimple
├── main.go                         # 约 10 行，调用 app.Run()
├── modules/
│   ├── demo/
│   │   ├── model.go                # GORM 模型
│   │   ├── controller.go           # 自定义控制器
│   │   ├── service.go              # 业务逻辑
│   │   └── module.go               # init() 注册入口
│   └── modules.go                  # 聚合 import 所有子模块
└── go.sum
```

### 快速上手

**1. 初始化项目并引入框架**

```bash
go mod init myapp
go get github.com/wilder2000/GOSimple
```

**2. 定义入口** `main.go`

```go
package main

import (
    "github.com/wilder2000/GOSimple/app"
    _ "myapp/modules"          // 业务模块自注册
)

func main() {
    app.Run()                  // 框架提供完整启动逻辑
}
```

**3. 创建业务模块** `modules/demo/model.go`

```go
package demo

type DemoOrder struct {
    ID     int32   `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID string  `gorm:"size:64;index" json:"user_id"`
    Amount float64 `json:"amount"`
    Status int32   `json:"status"`
}

func (DemoOrder) TableName() string { return "demo_order" }
```

**4. 注册通用 CRUD + 自定义接口** `modules/demo/module.go`

```go
package demo

import (
    "github.com/wilder2000/GOSimple/database"
    "github.com/wilder2000/GOSimple/http"
)

func init() {
    // 注册通用 CRUD —— 零代码生成增删改查接口
    http.RegObject[DemoOrder]("demo_order")

    // 注册自定义业务控制器
    http.RegMapping[OrderCreateRequest](&OrderCreateController{})

    // 幂等建表（仅新增，不破坏已有数据）
    database.DBHander.AutoMigrate(&DemoOrder{})
}
```

**5. 定义业务控制器** `modules/demo/controller.go`

```go
package demo

import (
    "github.com/gin-gonic/gin"
    "github.com/wilder2000/GOSimple/comm"
    "github.com/wilder2000/GOSimple/http"
)

type OrderCreateRequest struct {
    UserID string  `json:"user_id" validate:"required"`
    Amount float64 `json:"amount" validate:"required,gt=0"`
}

type OrderCreateController struct {
    http.AbstractController[OrderCreateRequest]
}

func (c *OrderCreateController) UrlPath() string { return "/demo/order/create" }

func (c *OrderCreateController) Execute(para *OrderCreateRequest, ctx *gin.Context) {
    // 业务逻辑...
    comm.SuccessResponse(ctx, gin.H{"order_id": 1})
}

// 实现 OperatorProvider 以纳入框架 RBAC 权限体系（可选，默认归管理员）
func (c *OrderCreateController) OperatorId() int32 { return 11 }
```

**6. 聚合 import** `modules/modules.go`

```go
package modules

import (
    _ "myapp/modules/demo"
    _ "myapp/modules/payment"
)
```

### 启动

```bash
go run .
go run . -install YES    # 初始化数据库
```

### 初始化执行顺序

```
config.init() → database.init() → glog.init() → http.init() → modules/*.init()
```

- `http.RegObject` / `http.RegMapping` 所需映射表在 `http.init()` 中已初始化
- `database.DBHander` 在 `database.init()` 中已可用
- 业务模块的 `init()` 最后执行，可安全调用上述 API

### 自动建表策略

二开模块的表由 `init()` 中的 `AutoMigrate` 自动创建，**早于 `main()` 执行**，不受启动模式影响：

```
执行时序：
  init() 链运行                    main() 运行
  │                              │
  ├─ config.init()               ├─ "-install YES" → dbscript.Install()
  │                              │   └─ DROP/CREATE 框架核心表 + 种子数据 + 模块 SQL
  ├─ ...                         │
  ├─ modules/demo.init()         ├─ 正常启动 → hs.Start()
  │   └─ AutoMigrate(&DemoOrder) │
  │       ↑ 自定义表已在此创建      │
```

无论 `go run .`（正常启动）还是 `go run . -install YES`，模块 `init()` 都会执行，`AutoMigrate` 都会运行。差异在于：
- **正常启动**：`init()` 建表 → `Start()` → HTTP 服务就绪
- **`-install YES`**：`init()` 建表 → `Install()` 重置框架核心表 + 种子数据 + 模块 SQL → 退出

| 阶段 | 推荐方式 | 说明 |
|------|----------|------|
| 开发期 | `AutoMigrate` 写在 `init()` 中 | 幂等，随服务启动自动同步 |
| 生产环境 | 独立迁移脚本 | 固化 DDL 变更，版本可控 |

`AutoMigrate` 行为：仅新增表和列，不修改/删除已有字段，**数据安全**。

### 模块安装 SQL

如果模块需要更精细的数据库初始化（如表结构变更、种子数据），可以通过 `dbscript.RegisterInstallSQL()` 注册安装 SQL，在 `-install YES` 时自动执行。

**1. 创建 SQL 文件** `modules/demo/install.sql`

```sql
DROP TABLE IF EXISTS demo_order;
CREATE TABLE demo_order (
    id       INT AUTO_INCREMENT PRIMARY KEY,
    user_id  VARCHAR(64) NOT NULL,
    amount   DECIMAL(10,2) NOT NULL,
    status   INT DEFAULT 0
);

INSERT INTO demo_order (user_id, amount, status) VALUES ('seed_user', 100.00, 1);
```

**2. 注册到安装流程** `modules/demo/module.go`

```go
package demo

import (
    _ "embed"
    "github.com/wilder2000/GOSimple/database"
    "github.com/wilder2000/GOSimple/dbscript"
    "github.com/wilder2000/GOSimple/http"
)

//go:embed install.sql
var installSQL string

func init() {
    http.RegObject[DemoOrder]("demo_order")
    http.RegMapping[OrderCreateRequest](&OrderCreateController{})
    database.DBHander.AutoMigrate(&DemoOrder{})

    // 注册安装 SQL，-install YES 时在框架核心表之后自动执行
    dbscript.RegisterInstallSQL("demo", installSQL)
}
```

**3. 执行安装**

```bash
go run . -install YES
```

**安装完整顺序**：
```
同一事务内:
  step 1  框架核心表 (initdb.sql)
  step 2  URL 映射
  step 3  种子数据 (管理员、角色、组)
  step 4  模块注册的 install SQL ← 新增，按注册顺序执行
```

**与 AutoMigrate 的关系**：
- `init()` 阶段：`AutoMigrate` 先执行（幂等建表，正常启动时处理演化）
- `-install YES` 阶段：模块 SQL 后执行（`DROP TABLE IF EXISTS` + `CREATE TABLE` 重建，确保表结构与 SQL 完全一致）

### 二开可复用的框架能力

| 能力 | 入口 | 说明 |
|------|------|------|
| 通用 CRUD | `http.RegObject[T]("name")` | 一行注册，自动生成 `/mif/{c,q,u,d}` 接口 |
| 自定义接口 | `http.RegMapping[T](ctrl)` | 实现 `HTTPController[T]` 接口即可 |
| 免认证接口 | `http.RegNoAuthMapping(path, handler)` | 注册不经过 JWT+RBAC 的路由，适合回调/webhook |
| RBAC 权限 | `OperatorId() int32` | 为 Controller 声明权限归属 |
| 线程池 | `pool.New[T](size, th, proc)` | 泛型异步任务处理 |
| 日志 | `glog.Logger.InfoF(...)` | 基于 Zap + Lumberjack |
| Excel 导出 | 通用查询 `attach: true` | 自动导出为 `.xlsx` |
| 验证器 | `go-playground/validator` | 内置中文翻译 |
| UUID / MD5 / bcrypt | `comm.UUID()` / `comm.MD5()` / `comm.EPassword()` | 开箱即用 |
| 模块安装 SQL | `dbscript.RegisterInstallSQL(name, sql)` | embed SQL 文件，`-install YES` 时自动执行 |

---

## 项目结构

```
├── main.go                   # 入口文件（薄封装，调用 app.Run()）
├── adminui.go                # 嵌入管理后台 SPA (//go:embed)
├── app/
│   └── app.go                # Run() 启动函数（可被外部项目导入）
├── conf/
│   ├── Application.yaml      # 应用配置
│   └── log4g.yaml            # 日志配置
├── http/
│   ├── http-server.go        # HTTP 服务初始化与启动
│   ├── http-auth.go          # JWT 令牌生成与验证
│   ├── http-dispatcher.go    # 前置中间件（认证 + 鉴权）
│   ├── http-command.go       # 响应模型与控制器抽象
│   ├── http-query.go         # 分页查询与 SQL 构建
│   ├── http-code.go          # 错误码定义
│   ├── mif-initial.go        # MIF 通用 CRUD 注册
│   ├── mif-model.go          # MIF 请求/响应模型
│   ├── mif-controller.go     # 通用 CRUD 实现
│   ├── handle-*.go           # HTTP 请求分发
│   ├── sm-*.go               # 用户/角色/组/部门管理
│   ├── url-sync.go           # URL 注册表与自动同步
│   └── validator-config.go   # 验证器 + 中文翻译
├── dbmodel/                  # GORM 数据模型
├── dbscript/                 # 数据库安装脚本
├── database/                 # 数据库连接与 SQL 构建器
├── config/                   # 配置管理
├── glog/                     # 日志封装
├── pool/                     # 泛型线程池
├── comm/                     # 通用工具函数
├── service/                  # 服务层
└── web/                      # 管理后台前端 (Vue 3 + Bootstrap)
    ├── src/
    │   ├── api/              # API 调用封装
    │   ├── stores/           # Pinia 状态管理
    │   ├── router/           # 路由配置
    │   ├── components/       # 通用组件
    │   └── views/            # 页面视图
    ├── index.html
    ├── vite.config.ts
    └── package.json
```
