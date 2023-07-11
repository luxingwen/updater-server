package model

import "gorm.io/gorm"

func MigrateDbTable(db *gorm.DB) {
	db.AutoMigrate(
		&DangerousCommand{},
		&FileInfo{},
		&PreTask{},
		&ScriptLibrary{},
	)

}
