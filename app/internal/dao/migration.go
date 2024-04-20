package dao

import (
	g "main/app/global"
	"main/app/internal/model"
)

// Migration 执行数据迁移
func Migration() {
	// 自动迁移模式
	err := g.MysqlDB.Set("gorm:table_options", "CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci").
		AutoMigrate(&model.UserSubject{}, &model.UserCollection{})
	if err != nil {
		return
	}
}
