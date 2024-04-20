package restaurant

import (
	"github.com/gin-gonic/gin"
	"main/app/internal/model"
	"main/app/internal/service"
	"net/http"
)

// Api 定义一个API的结构体
type Api struct{}

// Search 从Yelp搜索餐厅
func (a *Api) Search(c *gin.Context) {
	// 从请求中获取纬度、经度和搜索词
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	term := c.Query("term")

	// 从Yelp搜索餐厅
	res, err := service.Restaurant().Info().SearchRestaurantFromYelp(latitude, longitude, term)
	// 如果搜索过程中出现错误，返回错误
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

	// 创建一个餐厅的列表
	restaurants := new([]*model.Restaurant)
	// 将餐厅的数据解析为餐厅的对象
	err = service.Restaurant().Info().UnmarshalRestaurantsData(res, restaurants)
	// 如果解析过程中出现错误，返回错误
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

	// 返回成功响应，包括餐厅的列表
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "get restaurant successfully",
		"ok":   true,
		"data": restaurants,
	})
}
