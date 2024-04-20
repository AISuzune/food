package recipe

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	g "main/app/global"
	"main/app/internal/model"
	"main/app/internal/service"
	"net/http"
)

// Api 定义一个API的结构体
type Api struct{}

// Search 搜索菜谱
func (a *Api) Search(c *gin.Context) {
	// 从请求中获取饮食习惯、烹饪时间、准备时间、总时间、口味和食材
	dietary := c.Query("dietary")
	cookTimeString := c.Query("cook_time")
	perpTimeString := c.Query("perp_time")
	totalTimeString := c.Query("total_time")
	taste := c.QueryArray("taste")
	ingredients := c.QueryArray("ingredients")

	// 如果饮食习惯不为空并且不是halal、vegan或vegetarian，返回错误
	if dietary != "" {
		if dietary != "halal" && dietary != "vegan" && dietary != "vegetarian" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "invalid dietary",
				"ok":   false,
			})
			return
		}
	}

	// 获取烹饪时间、准备时间和总时间的时间段
	cookBeginTime, cookEndTime := service.Recipe().Info().GetTimeDuration(cookTimeString)
	perpBeginTime, perpEndTime := service.Recipe().Info().GetTimeDuration(perpTimeString)
	totalBeginTime, totalEndTime := service.Recipe().Info().GetTimeDuration(totalTimeString)

	// 从请求中获取限制数和页数，并将它们转换为整数
	limit := cast.ToInt64(c.Query("limit"))
	page := cast.ToInt64(c.Query("page"))

	// 如果限制数小于等于0，返回错误
	if limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  `invalid param "limit"`,
			"ok":   false,
		})
		return
	}
	// 如果页数小于等于0，返回错误
	if page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  `invalid param "page"`,
			"ok":   false,
		})
		return
	}

	// 获取数据库
	database := g.MongoDB.Database("food")

	// 获取菜谱的集合
	collection := database.Collection("recipe")

	// 定义一个过滤器
	filter := bson.D{}

	// 如果饮食习惯不为空
	if dietary != "" {
		// 如果饮食习惯是halal
		if dietary == "halal" {
			// 在过滤器中添加饮食习惯的条件，排除非halal的菜谱
			filter = append(filter, bson.E{
				Key: "dietary",
				Value: bson.D{
					{
						"$not",
						bson.D{
							{
								"$in",
								[]string{"non-halal"},
							},
						},
					},
				},
			})
			// 如果饮食习惯是vegetarian
		} else if dietary == "vegetarian" {
			// 在过滤器中添加饮食习惯的条件，排除非vegetarian的菜谱
			filter = append(filter, bson.E{
				Key: "dietary",
				Value: bson.D{
					{
						"$nin",
						[]string{"non-vegetarian"},
					},
				},
			})
			// 如果饮食习惯是vegan
		} else if dietary == "vegan" {
			// 在过滤器中添加饮食习惯的条件，排除非vegan的菜谱
			filter = append(filter, bson.E{
				Key: "dietary",
				Value: bson.D{
					{
						"$nin",
						[]string{"non-vegan"},
					},
				},
			})
		}
	}

	if cookEndTime != 0 {
		// 如果烹饪结束时间不为0，将烹饪时间添加到过滤器中
		filter = append(filter, bson.E{
			Key: "cook_time",
			Value: bson.D{
				{
					"$gte",
					cookBeginTime.Seconds(),
				},
				{
					"$lte",
					cookEndTime.Seconds(),
				},
			},
		})
	}

	if perpEndTime != 0 {
		// 如果准备结束时间不为0，将准备时间添加到过滤器中
		filter = append(filter, bson.E{
			Key: "perp_time",
			Value: bson.D{
				{
					"$gte",
					perpBeginTime.Seconds(),
				},
				{
					"$lte",
					perpEndTime.Seconds(),
				},
			},
		})
	}

	if totalEndTime != 0 {
		// 如果总结束时间不为0，将总时间添加到过滤器中
		filter = append(filter, bson.E{
			Key: "total_time",
			Value: bson.D{
				{
					"$gte",
					totalBeginTime.Seconds(),
				},
				{
					"$lte",
					totalEndTime.Seconds(),
				},
			},
		})
	}

	for _, ingredient := range ingredients {
		// 对于食材列表中的每一个食材，将其添加到过滤器中
		reg := bson.E{
			Key: "ingredients",
			Value: primitive.Regex{
				Pattern: ingredient, //匹配正则表达式ingredient的文档
				Options: "i",        //忽略大小写
			},
		}
		filter = append(filter, reg)
	}

	for _, tas := range taste {
		// 对于口味列表中的每一个口味，将其添加到过滤器中
		reg := bson.E{
			Key: "keywords",
			Value: primitive.Regex{
				Pattern: tas,
				Options: "i",
			},
		}
		filter = append(filter, reg)
	}

	// 打印过滤器的内容
	g.Logger.Debugf("%v", filter)

	// 创建一个查找选项，并设置限制数和跳过的数量
	option := &options.FindOptions{}
	option.SetLimit(limit)
	option.SetSkip(limit * (page - 1))
	// 在集合中查找匹配的文档
	cur, err := collection.Find(c, filter, option)
	if err != nil {
		// 如果查找过程中出现错误，记录错误日志并返回内部错误
		g.Logger.Errorf("find [recipe] document failed, err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "internal err",
			"ok":   false,
		})
		return
	}
	// 在函数返回后关闭游标
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			g.Logger.Errorf("close [recipe] document failed, err: %v", err)
		}
	}(cur, c)

	// 创建一个菜谱的列表
	var results []*model.Recipe

	// 遍历游标中的每一个文档
	for cur.Next(c) {
		// 创建一个菜谱的对象
		var elem model.Recipe

		// 将文档解码为菜谱的对象
		err := cur.Decode(&elem)
		if err != nil {
			// 如果解码过程中出现错误，记录错误日志并返回内部错误
			g.Logger.Errorf("decode [recipe] document failed, err: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "internal err",
				"ok":   false,
			})
			return
		}

		// 将菜谱添加到菜谱的列表中
		results = append(results, &elem)
	}

	// 返回成功响应，包括菜谱的列表
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "get recipe successfully",
		"ok":   true,
		"data": results,
	})
}
