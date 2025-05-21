package database

import (
	"fmt"

	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/glog"
	"gorm.io/gorm"
)

var (
	DBHander = &gorm.DB{}
)

func init() {
	LoadDatabaseConfig()
}
func LoadDatabaseConfig() {
	glog.Logger.Info("database begin to init.")
	if config.AConfig.DataSource.IsMySQL() {
		DBHander = createMySQLDbHandler()
	} else if config.AConfig.DataSource.IsSQLLite() {
		DBHander = createSQLiteDbHandler()
	} else {
		panic(fmt.Sprintf("Not supported database type %s", config.AConfig.DataSource.Type))
	}

	glog.Logger.Info("database init success.")
	printDb(config.AConfig.DataSource)
}
func printDb(db config.DBConfig) {
	glog.Logger.DebugF("Init Database%s", db.Name)
	glog.Logger.DebugF(db.Name+" %s", db.Type)
	glog.Logger.DebugF(db.Name+" %s", db.DSN)
	glog.Logger.DebugF(db.Name+"max Idle connections=%d", db.MaxIdleConnections)
	glog.Logger.DebugF(db.Name+"max Open connections=%d", db.MaxOpenConnections)
}

func Like(para string) string {
	return "%" + para + "%"
}
