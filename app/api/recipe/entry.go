package recipe

type Group struct{}

// insRecipe 创建一个API的实例
var insRecipe = Api{}

func (g *Group) Recipe() *Api {
	return &insRecipe
}
