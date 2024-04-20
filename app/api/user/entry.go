package user

type Group struct{}

// insSign 创建一个登录API的实例
var insSign = SignApi{}

func (g *Group) Sign() *SignApi {
	return &insSign
}

// insCollect 创建一个收藏API的实例
var insCollect = CollectApi{}

func (g *Group) Collect() *CollectApi {
	return &insCollect
}
