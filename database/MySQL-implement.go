package database

import (
	"fmt"

	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblog "gorm.io/gorm/logger"
)

func createMySQLDbHandler() *gorm.DB {
	dC := config.AConfig.DataSource
	myConfig := mysql.Config{
		DSN:               dC.DSN,
		DefaultStringSize: 256,
	}

	sqlLog := dblog.Default
	switch glog.LConfig.LogLevel {
	case glog.LevelDebug:
		sqlLog.LogMode(dblog.Info)
		fmt.Println("database log leve: Info")
	case glog.LevelInfo:
		sqlLog.LogMode(dblog.Info)
		fmt.Println("database log leve: info")
	case glog.LevelError:
		sqlLog.LogMode(dblog.Error)
		fmt.Println("database log leve: error")
	default:
		sqlLog.LogMode(dblog.Silent)
		fmt.Println("database log leve: default silent")
	}
	glog.Logger.InfoF("DSN:%s", myConfig.DSN)
	db, err := gorm.Open(mysql.New(myConfig), &gorm.Config{
		Logger: sqlLog,
	})
	glog.Logger.Info("database opened.")
	if err != nil {
		panic(err)
	}
	sqldb, err2 := db.DB()
	if err2 != nil {
		glog.Logger.ErrorF("Create DB Pool failed %s", err2.Error())
	}

	sqldb.SetMaxOpenConns(dC.MaxOpenConnections)
	sqldb.SetMaxIdleConns(dC.MaxIdleConnections)
	glog.Logger.Info("database pool init.")
	return db
}
