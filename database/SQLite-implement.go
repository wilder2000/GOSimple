package database

import (
	"sync"
	"time"

	"github.com/wilder2000/GOSimple/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	dblog "gorm.io/gorm/logger"
)

// import (
//
//	"gorm.io/gorm"
//	"wilder.cn/gogo/log"
//
// )
// import "gorm.io/driver/sqlite"
//
//	type SQLiteHandler struct {
//		DbFile string
//	}
//
//	func (r SQLiteHandler) Open() (*gorm.DB, bool) {
//		db, err := gorm.Open(sqlite.Open(r.DbFile), &gorm.Config{})
//		if err != nil {
//			log.Logger.ErrorF("Try to open sqlite failed.%s", err.Error())
//
//			return nil, false
//		}
//		return db, true
//	}

func createSQLiteDbHandler() *gorm.DB {
	dbConfig := config.AConfig.DataSource

	var err error
	db, err := gorm.Open(sqlite.Open(dbConfig.DSN), &gorm.Config{
		Logger: dblog.Default.LogMode(dblog.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 获取底层sql.DB连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 启用SQLite WAL模式提升并发性能
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	return db

}
