package user

import (
	"context"
	g "main/app/global"
	"main/app/internal/model"
)

// DUser 是一个结构体，用于处理用户相关的操作
type DUser struct{}

func (d *DUser) GetUserByUsername(ctx context.Context, username string) error {
	// 创建一个用户对象
	userSubject := &model.UserSubject{}
	// 在数据库中查找用户名
	err := g.MysqlDB.WithContext(ctx).
		Table("user_subject").
		Select("username").
		Where("username = ?", username).
		First(userSubject).Error
	return err
}

func (d *DUser) CreateUser(ctx context.Context, userSubject *model.UserSubject) {
	// 在数据库中创建用户
	g.MysqlDB.WithContext(ctx).
		Table("user_subject").
		Create(userSubject)
}

func (d *DUser) GetUserByUsernameAndPassword(ctx context.Context, userSubject *model.UserSubject) error {
	// 在数据库中查找用户名和密码都匹配的用户
	err := g.MysqlDB.WithContext(ctx).
		Table("user_subject").
		Where(&model.UserSubject{
			Username: userSubject.Username,
			Password: userSubject.Password,
		}).
		First(userSubject).Error
	return err
}
