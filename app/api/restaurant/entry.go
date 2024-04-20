package restaurant

type Group struct{}

// insRestaurant 创建一个API的实例
var insRestaurant = Api{}

func (g *Group) Restaurant() *Api {
	return &insRestaurant
}
