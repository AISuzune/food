package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	g "main/app/global"
	"main/app/internal/model"
	"main/app/internal/service"
	"net/http"
)

// CollectApi 定义一个收藏API的结构体
type CollectApi struct{}

// GetList 获取收藏列表的函数
func (a *CollectApi) GetList(c *gin.Context) {
	// 从请求中获取用户ID
	userId := c.GetInt64("id")

	// 从请求中获取收藏类型、限制数和页数，并将它们转换为整数
	collectType := cast.ToInt32(c.Query("collect_type"))
	limit := cast.ToInt32(c.Query("limit"))
	page := cast.ToInt32(c.Query("page"))

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

	// 定义一个用户收藏的列表
	var userCollections []*model.UserCollection
	// 根据收藏类型来获取不同的收藏列表
	switch collectType {
	case 1, 2:
		// 获取用户的收藏数量
		cnt, err := service.User().Collect().GetUserCollectionCount(c, userId, collectType)
		// 如果(收藏数量为-1)或者出现错误，返回错误
		if err != nil {
			switch err.Error() {
			case "internal err":
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "internal err",
					"ok":   false,
				})
			}

			return
		}

		// 计算页数
		pageCount := int32(cnt) / limit
		if int32(cnt)%limit > 0 {
			pageCount = pageCount + 1
		}

		// 如果页数大于最大页数，返回错误
		if page > pageCount {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  fmt.Sprintf("the maximum number of pages is %d", pageCount),
			})
			return
		}

		// 获取用户的收藏列表
		userCollections, err = service.User().Collect().GetUserCollectionsWithLimit(c, userId, collectType, int(limit), int(page))
		// 如果出现错误，返回错误
		if err != nil {
			switch err.Error() {
			case "internal err":
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "internal err",
					"ok":   false,
				})
			}

			return
		}

	// 如果收藏类型无效，返回错误
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid collect type",
			"ok":   false,
		})
		return
	}

	// 根据收藏类型来返回不同的收藏列表
	switch collectType {
	case 1:
		// 定义一个收藏的列表
		var collections []*model.Collection
		// 遍历用户收藏的列表
		for _, userCollection := range userCollections {
			// 获取餐厅的信息
			res, err := service.Restaurant().Info().GetRestaurantByID(userCollection.RestaurantId)
			// 如果出现错误，返回错误
			if err != nil {
				switch err.Error() {
				case "internal err":
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": http.StatusInternalServerError,
						"msg":  "internal err",
						"ok":   false,
					})
				}

				return
			}

			// 创建一个餐厅的对象
			restaurant := &model.Restaurant{}
			// 将餐厅的信息解析为JSON
			err = json.Unmarshal([]byte(res), restaurant)
			// 如果解析过程中出现错误，返回错误
			if err != nil {
				g.Logger.Errorf("unmarshal restaurant json failed, err: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "internal err",
					"ok":   false,
				})
			}
			// 创建一个收藏的对象，并设置收藏的类型和数据
			collection := &model.Collection{
				Id:             userCollection.Id,
				CollectionType: "restaurant",
				CollectionData: restaurant,
			}
			// 将收藏添加到收藏的列表中
			collections = append(collections, collection)
		}

		// 返回成功的响应
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "get collection successfully",
			"ok":   true,
			"data": collections,
		})

	case 2:
		// 定义一个收藏的列表
		var collections []*model.Collection
		// 遍历用户收藏的列表
		for _, userCollection := range userCollections {
			// 获取菜谱的信息
			recipe := service.Recipe().Info().GetRecipeById(c, userCollection.RecipeId)
			// 如果菜谱存在
			if recipe != nil {
				// 创建一个收藏的对象，并设置收藏的类型和数据
				collection := &model.Collection{
					Id:             userCollection.Id,
					CollectionType: "recipe",
					CollectionData: recipe,
				}
				// 将收藏添加到收藏的列表中
				collections = append(collections, collection)
			}
		}

		// 返回成功的响应
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "get collection successfully",
			"ok":   true,
			"data": collections,
		})
	}
}

func (a *CollectApi) Create(c *gin.Context) {
	// 从上下文中获取用户ID
	userId := c.GetInt64("id")

	// 从表单中获取收藏类型，并将其转换为整数
	collectType := cast.ToInt32(c.PostForm("collect_type"))

	// 创建一个用户收藏的对象，并设置用户ID和收藏类型
	userCollection := &model.UserCollection{
		UserId:      userId,
		CollectType: collectType,
	}

	// 定义一个空的ID
	var id interface{}

	// 根据收藏类型来处理不同的收藏
	switch collectType {
	case 1:
		// 收藏餐厅
		restaurantId := c.PostForm("restaurant_id")
		// 如果餐厅ID为空，返回错误
		if restaurantId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "restaurant_id cannot be null",
				"ok":   false,
			})
			return
		}

		// 将餐厅ID赋值给ID
		id = restaurantId
		// 设置用户收藏的餐厅ID
		userCollection.RestaurantId = restaurantId

	case 2:
		// 收藏菜谱
		recipeId := cast.ToInt64(c.PostForm("recipe_id"))
		// 如果菜谱ID为0，返回错误
		if recipeId == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "recipe_id cannot be null",
				"ok":   false,
			})
			return
		}

		// 将菜谱ID赋值给ID
		id = recipeId
		// 设置用户收藏的菜谱ID
		userCollection.RecipeId = recipeId

	default:
		// 如果收藏类型无效，返回错误
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid collect type",
			"ok":   false,
		})
		return
	}

	// 检查收藏是否已存在
	err := service.User().Collect().CheckCollectionIsExist(c, collectType, userId, id)
	if err != nil {
		switch err.Error() {
		case "internal err":
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "internal err",
				"ok":   false,
			})

		case "duplicate collect":
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"ok":   false,
			})
		}

		return
	}

	// 在数据库中创建收藏
	service.User().Collect().CreateCollection(c, userCollection)
	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "create collection successfully",
		"ok":   true,
	})
}

func (a *CollectApi) Delete(c *gin.Context) {
	// 从请求中获取ID，并将其转换为整数
	id := cast.ToInt64(c.Query("id"))
	// 从上下文中获取用户ID
	userId := c.GetInt64("id")

	// 如果ID为0，返回错误
	if id == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "id cannot be null",
			"ok":   false,
		})
		return
	}

	// 检查收藏ID是否存在
	err := service.User().Collect().CheckCollectionIdIsExist(c, id, userId)
	if err != nil {
		switch err.Error() {
		case "internal err":
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "internal err",
				"ok":   false,
			})
		case "collection not found":
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"msg":  "collection not found",
				"ok":   false,
			})
		}

		return
	}

	// 在数据库中删除收藏
	err = service.User().Collect().DeleteCollection(c, id)
	if err != nil {
		switch err.Error() {
		case "internal err":
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "internal err",
				"ok":   false,
			})
		}
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "delete collection successfully",
		"ok":   true,
	})
}
