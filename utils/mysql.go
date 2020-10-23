package utils

import "gorm.io/gorm"

func CreateTable(db *gorm.DB, tableName string, tableStruct interface{}) error {
	migrator := db.Table(tableName).Migrator()
	if !migrator.HasTable(tableName) {
		if err := migrator.CreateTable(tableStruct); err != nil {
			return err
		}
	}
	return nil
}
