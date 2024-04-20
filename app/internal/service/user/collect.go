package user

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	g "main/app/global"
	"main/app/internal/dao"
	"main/app/internal/model"
)

// SCollect 定义一个收藏的结构体，，用于处理收藏相关的操作
type SCollect struct{}

// CheckCollectionIsExist 检查给定的收藏是否存在
func (s *SCollect) CheckCollectionIsExist(ctx context.Context, collectType int32, userId int64, id interface{}) error {
	// 在数据库中查找是否存在这样的收藏
	err := dao.User().Collect().GetCollectionByType(ctx, collectType, userId, id)
	// 如果查找过程中出现错误
	if err != nil {
		// 如果错误不是因为找不到记录
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录错误日志
			g.Logger.Errorf("query [user_collection] record failed, err: %v", err)
			// 返回内部错误
			return fmt.Errorf("internal err")
		}
		// 如果找到了收藏
	} else {
		// 返回收藏已存在的错误
		return fmt.Errorf("duplicate collect")
	}

	// 如果找不到记录，返回nil
	return nil
}

// CheckCollectionIdIsExist 检查给定的收藏ID是否存在
func (s *SCollect) CheckCollectionIdIsExist(ctx context.Context, id, userId int64) error {
	// 在数据库中查找是否存在一个ID和用户ID都匹配的收藏
	err := dao.User().Collect().GetCollectionById(ctx, id, userId)
	// 如果查找过程中出现错误
	if err != nil {
		// 如果错误是因为找不到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回收藏不存在的错误
			return fmt.Errorf("collection not found")
		}
		// 记录错误日志
		g.Logger.Errorf("query [user_collection] record failed, err: %v", err)
		// 返回内部错误
		return fmt.Errorf("internal err")
	}

	// 如果没有错误，返回nil
	return nil
}

// CreateCollection 在数据库中创建一个新的收藏
func (s *SCollect) CreateCollection(ctx context.Context, userCollection *model.UserCollection) {
	// 在数据库中创建收藏
	dao.User().Collect().CreateCollection(ctx, userCollection)
}

// DeleteCollection 在数据库中删除一个收藏
func (s *SCollect) DeleteCollection(ctx context.Context, id int64) error {
	// 在数据库中删除这个ID对应的收藏
	err := dao.User().Collect().DeleteCollection(ctx, id)
	// 如果删除过程中出现错误
	if err != nil {
		// 记录错误日志
		g.Logger.Errorf("delete [user_collection] record failed, err: %v", err)
		// 返回内部错误
		return fmt.Errorf("internal err")
	}

	// 如果没有错误，返回nil
	return nil
}

// GetUserCollectionCount 获取用户的收藏数量
func (s *SCollect) GetUserCollectionCount(ctx context.Context, userId int64, collectType int32) (int64, error) {
	// 在数据库中计算这个用户的这种类型的收藏的数量
	cnt, err := dao.User().Collect().GetUserCollectionCount(ctx, userId, collectType)
	// 如果计算过程中出现错误
	if err != nil {
		// 记录错误日志
		g.Logger.Errorf("query [user_collection] record failed ,err: %v", err)
		// 返回-1和内部错误
		return -1, fmt.Errorf("internal err")
	}

	// 如果没有错误，返回收藏的数量
	return cnt, nil
}

// GetUserCollectionsWithLimit 获取用户的收藏列表
func (s *SCollect) GetUserCollectionsWithLimit(ctx context.Context, userId int64, collectType int32, limit, page int) ([]*model.UserCollection, error) {
	// 在数据库中分页查找这个用户的这种类型的收藏的列表
	userCollections, err := dao.User().Collect().GetUserCollectionsWithLimit(ctx, userId, collectType, limit, page)
	// 如果查找过程中出现错误
	if err != nil {
		// 记录错误日志
		g.Logger.Errorf("query [user_collection] failed, err: %v", err)
		// 返回nil和内部错误
		return nil, fmt.Errorf("internal err")
	}

	// 如果没有错误，返回收藏的列表
	return userCollections, nil
}
