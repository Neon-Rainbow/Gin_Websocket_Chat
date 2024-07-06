package MySQL

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"websocket/config"
	"websocket/model"
)

var MySQL *gorm.DB

func InitMySQL() (db *gorm.DB, err error) {
	// 初始化MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.MySQL.User,
		config.AppConfig.MySQL.Password,
		config.AppConfig.MySQL.Host,
		config.AppConfig.MySQL.Port,
		config.AppConfig.MySQL.DBName,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("服务器连接失败: %v", err)
		return nil, err
	}
	err = db.AutoMigrate(
		&model.SQLMessage{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
		return nil, err
	}

	MySQL = db

	return db, nil
}
