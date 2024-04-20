package recipe

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	g "main/app/global"
	"main/app/internal/model"
	"strings"
	"time"
)

// SInfo 定义一个菜谱信息的结构体
type SInfo struct{}

// GetRecipeById 根据ID从数据库中获取菜谱
func (s *SInfo) GetRecipeById(ctx context.Context, recipeId int64) *model.Recipe {
	// 定义一个过滤器，用于在数据库中查找匹配的菜谱
	filter := bson.D{
		{
			"recipe_id",
			recipeId,
		},
	}

	// 创建一个菜谱的对象
	var elem model.Recipe

	// 在数据库中查找匹配的菜谱
	cur := g.MongoDB.Database("food").Collection("recipe").
		FindOne(ctx, filter)

	// 将查找到的菜谱解码为菜谱的对象
	err := cur.Decode(&elem)
	// 如果解码过程中出现错误，返回nil
	if err != nil {
		return nil
	}

	// 返回菜谱的对象
	return &elem
}

// GetTimeDuration 获取时间段
func (s *SInfo) GetTimeDuration(timeStr string) (time.Duration, time.Duration) {
	// 如果时间字符串为空，返回0和0
	if timeStr == "" {
		return 0, 0
	}
	// 将时间字符串按照"-"分割为两部分
	output := strings.Split(timeStr, "-")
	// 如果分割后的长度不为2，返回0和0
	if len(output) != 2 {
		return 0, 0
	}
	// 将分割后的两部分分别解析为时间段
	from, _ := time.ParseDuration(output[0])
	to, _ := time.ParseDuration(output[1])
	// 返回两个时间段
	return from, to
}
