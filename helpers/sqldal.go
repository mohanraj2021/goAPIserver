package helpers

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var MariaDB *gorm.DB

var PtrToLogger *zap.Logger

func InitMariaDB() (f bool) {
	f = false
	defer func() {
		if r := recover(); r != nil {
			f = false
			return
		}
	}()

	db, err := gorm.Open(mysql.Open(MariaDBConnStr), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		PtrToLogger.Fatal("sgt_portal_controller", zap.String("message", fmt.Sprintf("InitMariaDB:gorm connection error: %s", err.Error())))
		return f
	}
	f = true
	MariaDB = db
	return f
}
