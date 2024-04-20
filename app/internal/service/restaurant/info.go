package restaurant

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	g "main/app/global"
	"main/app/internal/model"
	"net/http"
)

// SInfo 定义一个餐厅信息的结构体
type SInfo struct{}

// 定义Yelp的搜索API和ID API的URL
const (
	yelpSearchApi = "https://api.yelp.com/v3/businesses/search"
	yelpIdApi     = "https://api.yelp.com/v3/businesses/"
)

// SearchRestaurantFromYelp 从Yelp搜索餐厅
func (s *SInfo) SearchRestaurantFromYelp(latitude, longitude, term string) (string, error) {
	// 发送一个带有授权头部、纬度、经度和搜索词的GET请求到Yelp的搜索API
	res, err := req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", g.Config.YelpApiKey)).
		SetQueryParams(map[string]string{
			"latitude":  latitude,
			"longitude": longitude,
			"term":      term,
		}).Get(yelpSearchApi)
	// 如果请求过程中出现错误，记录错误日志并返回内部错误
	if err != nil {
		g.Logger.Errorf("query yelp api failed, err: %v", err)
		return "", fmt.Errorf("internal err")
	}
	// 如果响应的状态码不是200，记录错误日志并返回内部错误
	if res.StatusCode != http.StatusOK {
		g.Logger.Errorf("query yelp api failed, err: %v", res)
		return "", fmt.Errorf("internal err")
	}
	// 返回响应的字符串
	return res.String(), nil
}

// UnmarshalRestaurantsData 将餐厅数据解析为餐厅的对象
func (s *SInfo) UnmarshalRestaurantsData(src string, dst *[]*model.Restaurant) error {
	// 将源字符串解析为JSON
	resJson := gjson.Parse(src)
	// 从JSON中获取餐厅的数组，从resJson这个JSON对象中获取名为"businesses"的值
	resArray := resJson.Get("businesses").Array()

	// 定义一个餐厅的列表
	var dstData []*model.Restaurant

	// 遍历餐厅的数组
	for _, restaurantRes := range resArray {
		// 创建一个餐厅的对象
		restaurant := &model.Restaurant{}
		// 将餐厅的JSON解析为餐厅的对象
		err := json.Unmarshal([]byte(restaurantRes.Raw), restaurant)
		// 如果解析过程中出现错误，记录错误日志并返回内部错误
		if err != nil {
			g.Logger.Errorf("unmarshal restaurant json failed, err: %v", err)
			return fmt.Errorf("internal err")
		}
		// 将餐厅添加到餐厅的列表中
		dstData = append(dstData, restaurant)
	}

	// 将餐厅的列表赋值给目标
	*dst = dstData

	// 返回nil表示没有错误
	return nil
}

// GetRestaurantByID 根据ID从Yelp获取餐厅
func (s *SInfo) GetRestaurantByID(id string) (string, error) {
	// 发送一个带有授权头部的GET请求到Yelp的ID API
	res, err := req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", g.Config.YelpApiKey)).
		Get(fmt.Sprintf("%s%s", yelpIdApi, id))
	// 如果请求过程中出现错误，记录错误日志并返回内部错误
	if err != nil {
		g.Logger.Errorf("query yelp api failed, err: %v", err)
		return "", fmt.Errorf("internal err")
	}
	// 如果响应的状态码不是200，记录错误日志并返回内部错误
	if res.StatusCode != http.StatusOK {
		g.Logger.Errorf("query yelp api failed, err: %v", res)
		return "", fmt.Errorf("internal err")
	}

	// 返回响应的字符串
	return res.String(), nil
}
