# GOSimple

Simple Go framework / util library. Module: `github.com/wilder2000/GOSimple` (Go 1.23.4).

## Commands

- `go run .` — starts the HTTP server (+ admin UI at /admin)
- `go run . -install YES` — initializes DB schema + seed data
- `go build .` — builds the binary
- `cd web && npm install && npm run build` — rebuilds admin frontend

No test/lint/format scripts exist.

## Architecture

- **Entrypoint**: `main.go` — calls `app.Run()`.
- **App bootstrap**: `app/app.go` — `Run()` function, importable by external projects.
- **HTTP**: Gin-based server in `http/` package. JWT auth (`dgrijalva/jwt-go`). Controller registration uses `RegMapping[T](controller)` with reflection-based dispatch. Serves on `:8090` by default (`conf/Application.yaml`).
- **Config**: Viper loads `conf/Application.yaml` + `conf/log4g.yaml`. Env overrides: `GOGO_HOME` (app dir), `GOGO_CONFIG_FILE` (config filename).
- **Database**: GORM with MySQL or SQLite (set via `DataSource.type` in config). DB init in `config.init()`, re-loadable via `database.LoadDatabaseConfig()`.
- **Logging**: `glog` package wraps Uber Zap + Lumberjack rotation.
- **Generic CRUD endpoints**: `/mif/c` (create), `/mif/q` (query), `/mif/u` (update), `/mif/d` (delete). Targets registered in `http/mif-initial.go` via `RegObject[T]("name")`.
- **Install**: `dbscript/install.go` embeds `dbscript/MySQL/initdb.sql` via `//go:embed`. Run with `-install YES` flag.
- **Thread pool**: Generic `pool.PoolEngine[T]` in `pool/poolengine.go`.
- **Models**: GORM structs in `dbmodel/` (`s_*.gen.go`), table names via `TableName()` method.

## Conventions

- Package naming: `dbmodel`, `dbscript`, `glog`, `comm` — not idiomatic Go but consistent internally.
- Files named `s_*.gen.go` are generated models; `sm-*.go` are user/security handlers.
- All response types are JSON with `{message, code, data}` format (see `http/http-command.go`).
- Validator translations default to Chinese (`zh`), registered in `http/validator-config.go`.
- Password hashing: `golang.org/x/crypto/bcrypt` via `comm.EPassword()`.
- SQL builder helper in `database/sqlbuilder.go` uses `goqu` with `mysql` dialect.
- `config.AConfig` is a global singleton populated at init time.
- `database.DBHander` is a global `*gorm.DB` populated at init time.
