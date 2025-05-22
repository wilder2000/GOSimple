package dbscript

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/glog"
	"gorm.io/gorm"
)

const (
	OPER_ID_ADMIN  = 10
	OPER_ID_VIEWER = 11

	SQL_1 = `INSERT INTO s_role VALUES (1,'终端用户角色','%s'),(2,'管理员角色','%s');` //exec must add time
	SQL_2 = `INSERT INTO s_roleoperator VALUES (1,11,0),(2,10,0),(2,11,0);`
	SQL_3 = `INSERT INTO s_rolegroup VALUES (1,1),(1,2),(2,2);`
	SQL_4 = `insert into s_group(id,name,createtime) value(1,'终端用户组','%s');` //exec must add time
	SQL_5 = `insert into s_group(id,name,createtime) value(2,'管理员组','%s');`  //exec must add time
	SQL_6 = `INSERT INTO s_operators VALUES (10,'管理功能'),(11,'查看功能');`
	SQL_7 = `INSERT INTO s_groupuser VALUES (2,'Administrator');`
	//1s:password,2s：time,3s:time
	SQL_8 = `INSERT INTO s_users VALUES ('Administrator','wild.shang@163.com','%s','流星划过黑暗的夜空つ','',0,'%s','%s',2,'','管理员');`
)

// 解析SQL文件，返回SQL语句切片
func parseSQLFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var statements []string
	var currentStmt strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过注释和空行
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}

		currentStmt.WriteString(" " + line)

		// 以分号作为语句结束符
		if strings.HasSuffix(line, ";") {
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
		}
	}

	return statements, scanner.Err()
}

func Install(authMappings []string) {
	db := database.DBHander

	homedir := config.AppDir()

	sqlFile := filepath.Join(homedir, "dbscript", "MySQL", "initdb.sql")
	statements, err := parseSQLFile(sqlFile)
	if err != nil {
		panic(err)
	}
	authSQLS := generateUrlMappingSQL(authMappings, OPER_ID_ADMIN)
	// authSQL := generateUrlMappingSQL(authMappings, OPER_ID_VIEWER)
	terr := db.Transaction(func(tx *gorm.DB) error {
		//step 1: execute create table sql.
		for _, stmt := range statements {
			if err := execOneSQL(tx, stmt); err != nil {
				return err
			}
		}
		//step 2: execute auth data.
		for _, oneSql := range authSQLS {
			if err := execOneSQL(tx, oneSql); err != nil {
				return err
			}
		}

		ntime := comm.NowTime()
		if err := execOneSQL(tx, fmt.Sprintf(SQL_1, ntime, ntime)); err != nil {
			return err
		}
		if err := execOneSQL(tx, SQL_2); err != nil {
			return err
		}
		if err := execOneSQL(tx, SQL_3); err != nil {
			return err
		}
		if err := execOneSQL(tx, fmt.Sprintf(SQL_4, ntime)); err != nil {
			return err
		}
		if err := execOneSQL(tx, fmt.Sprintf(SQL_5, ntime)); err != nil {
			return err
		}
		if err := execOneSQL(tx, SQL_6); err != nil {
			return err
		}
		if err := execOneSQL(tx, SQL_7); err != nil {
			return err
		}
		hpwd, err := comm.EPassword(config.AConfig.Security.DefaultAdminPWD)
		if err != nil {
			return errors.New(fmt.Sprintf("Default admin pwd init failed:%s", err.Error()))
		}
		if err := execOneSQL(tx, fmt.Sprintf(SQL_8, string(hpwd), ntime, ntime)); err != nil {
			return err
		}
		return nil
	})

	if terr != nil {
		glog.Logger.ErrorF("INIT database failed. %s", terr.Error())
	} else {
		glog.Logger.InfoF("INIT database success.")
	}

}
func execOneSQL(db *gorm.DB, sql string) error {
	if err := db.Exec(sql).Error; err != nil {
		return errors.New(fmt.Sprintf("sql exec error: %v\nSQL: %s\n", err, sql))
	} else {
		glog.Logger.InfoF("SUCCESS: %s", sql)
	}
	return nil
}
func generateUrlMappingSQL(urls []string, operatorID int) []string {
	builder := make([]string, 0)
	for _, url := range urls {
		one := fmt.Sprintf("insert into s_urlmappings(operatorid,url) values(%d,'%s');\n", operatorID, url)
		builder = append(builder, one)
	}
	return builder
}
