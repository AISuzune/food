package user

import (
	"context"
	"fmt"
	g "main/app/global"
	"main/app/internal/model"
)

// DCollect 定义一个收藏的结构体，，用于处理收藏相关的操作
type DCollect struct{}

func (d *DCollect) GetCollectionByType(ctx context.Context, collectType int32, userId int64, id interface{}) error {
	// 定义一个空的SQL语句
	whereSql := ""
	// 根据收藏类型来生成不同的SQL语句
	switch collectType {
	case 1:
		whereSql = fmt.Sprintf("user_id = ? AND restaurant_id = ?")
	case 2:
		whereSql = fmt.Sprintf("user_id = ? AND recipe_id = ?")
	default:

	}

	// 创建一个用户收藏的对象
	userCollection := &model.UserCollection{}
	// 在数据库中查找是否存在这样的收藏
	err := g.MysqlDB.WithContext(ctx).
		Table("user_collection").
		Where(whereSql, userId, id).
		First(userCollection).Error
	return err
}

func (d *DCollect) GetCollectionById(ctx context.Context, id, userId int64) error {
	// 创建一个用户收藏的对象
	userCollection := &model.UserCollection{}
	// 在数据库中查找是否存在一个ID和用户ID都匹配的收藏
	err := g.MysqlDB.WithContext(ctx).
		Table("user_collection").
		Select("id,user_id").
		Where("id= ? AND user_id = ?", id, userId).
		First(userCollection).Error
	return err
}

func (d *DCollect) CreateCollection(ctx context.Context, userCollection *model.UserCollection) {
	// 在数据库中创建收藏
	g.MysqlDB.WithContext(ctx).
		Table("user_collection").
		Create(userCollection)
}

func (d *DCollect) DeleteCollection(ctx context.Context, id int64) error {
	// 在数据库中删除这个ID对应的收藏
	err := g.MysqlDB.WithContext(ctx).
		Table("user_collection").
		Delete(&model.UserCollection{}, id).Error
	return err
}

func (d *DCollect) GetUserCollectionCount(ctx context.Context, userId int64, collectType int32) (int64, error) {
	// 定义一个计数器
	var cnt int64
	// 在数据库中计算这个用户的这种类型的收藏的数量
	err := g.MysqlDB.WithContext(ctx).
		Table("user_collection").
		Where("user_id = ? AND collect_type = ?", userId, collectType).
		Count(&cnt).Error
	return cnt, err
}

func (d *DCollect) GetUserCollectionsWithLimit(ctx context.Context, userId int64, collectType int32, limit, page int) ([]*model.UserCollection, error) {
	// 定义一个用户收藏的列表
	var userCollections []*model.UserCollection
	// 在数据库中分页查找这个用户的这种类型的收藏的列表
	err := g.MysqlDB.WithContext(ctx).
		Table("user_collection").
		Limit(limit).Offset(limit*(page-1)).
		Where("user_id = ? AND collect_type = ?", userId, collectType).
		Find(&userCollections).Error
	return userCollections, err
}
