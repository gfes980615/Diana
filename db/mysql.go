package db

import (
	"fmt"
	"github.com/gfes980615/Diana/glob/common/log"

	"github.com/gfes980615/Diana/glob/config"
	"gorm.io/driver/mysql" //加载mysql
	"gorm.io/gorm"
)

// MySQL ...
type MySQL struct {
	gorm *gorm.DB
}

var MysqlConn *MySQL

// InitMysql ...
func InitMysql(mysqlConf config.Mysql) error {

	MysqlConn = new(MySQL)
	if err := MysqlConn.Connect(mysqlConf); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// Connect ...
func (db *MySQL) Connect(mysqlConf config.Mysql) error {
	var err error
	config := mysqlConf.DataBases
	connect := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		config.Username,
		config.Password,
		config.Address,
		config.Database,
	)

	db.gorm, err = gorm.Open(mysql.Open(connect), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Debug("Database [%s] Connect success", config.Name)

	// TODO, load from config
	//db.gorm.LogMode(mysqlConf.LogMode)
	//
	//// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	//db.gorm.DB().SetMaxIdleConns(mysqlConf.MaxIdle)
	//// SetMaxOpenConns sets the maximum number of open connections to the database.
	//db.gorm.DB().SetMaxOpenConns(mysqlConf.MaxOpen)
	//// SetConnMaxLifetime sets the maximum amount of timeUtil a connection may be reused.
	//db.gorm.DB().SetConnMaxLifetime(time.Duration(mysqlConf.ConnMaxLifeMin) * time.Minute)

	return nil
}

// Session ...
func (db *MySQL) Session() *gorm.DB {
	return db.gorm
}

// Begin ...
func (db *MySQL) Begin() *gorm.DB {
	return db.gorm.Begin()
}

// Close ...
//func (db *MySQL) Close() {
//	db.gor
//}
