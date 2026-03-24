package database

import (
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 数据库初始化
func SetupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("network_security.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移数据库表
	db.AutoMigrate(
		&models.ProductType{},
		&models.Function{},
		&models.Product{},
		&models.NetworkTopo{},
	)

	// 初始化基础数据
	SeedData(db)
	return db
}
