package db

import (
	"fmt"
	"github.com/cihub/seelog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"time"
)

// InitDB : init db
func InitMysql(user, password, host, port, dbName string, maxIdls, maxConns, idlTimeout int) (*gorm.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&allowNativePasswords=true",
		user,
		password,
		host,
		port,
		dbName,
	)

	//DefaultLogger := logger.Default
	//if environment.Env.IsDev() || environment.Env.IsLocal() {
	//	DefaultLogger = DefaultLogger.LogMode(logger.Info)
	//}
	MysqlDB, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	//MysqlDB, err := gorm.Open("mysql", dataSource)
	if err != nil {
		seelog.Errorf("InitMysql dataSource=%+v, error=%s", dataSource, err.Error())
		return nil, err
	}
	// 初始化连接池
	sqlDB, err := MysqlDB.DB()
	if err != nil {
		seelog.Errorf("obtain DB(), error=%s", err.Error())
		return nil, err
	}
	sqlDB.SetMaxIdleConns(maxIdls)
	sqlDB.SetMaxOpenConns(maxConns)
	idleTimeout := idlTimeout
	if idleTimeout > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(idleTimeout) * time.Second)
	}

	MysqlDB.Config.Logger.LogMode(logger.Info)
	return MysqlDB, nil
}
