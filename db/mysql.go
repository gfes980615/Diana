package db

import (
	"fmt"
	"log"

	"github.com/gfes980615/Diana/model"
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type MySQL struct {
	DB *gorm.DB
}

// NewMySQL ...
func NewMySQL(config model.DataBaseConfig) (*MySQL, error) {
	db := new(MySQL)
	err := db.Connect(config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Connect ...
func (db *MySQL) Connect(config model.DataBaseConfig) error {
	var err error
	connect := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		glob.DataBase.Username,
		glob.DataBase.Password,
		glob.DataBase.Address,
		glob.DataBase.Database,
	)

	db.DB, err = gorm.Open("mysql", connect)
	if err != nil {
		return err
	}

	log.Printf("Database [%s] Connect success", config.Database)

	return nil
}

// Session ...
func (db *MySQL) Session() *gorm.DB {
	return db.DB
}

// Begin ...
func (db *MySQL) Begin() *gorm.DB {
	return db.DB.Begin()
}

// Close ...
func (db *MySQL) Close() {
	db.DB.Close()
}
