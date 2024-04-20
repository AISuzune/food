package user

type Group struct{}

// insUser 是 DUser 的一个实例
var insUser = DUser{}

func (g *Group) User() *DUser {
	return &insUser
}

// insCollect 创建一个收藏的实例
var insCollect = DCollect{}

func (g *Group) Collect() *DCollect {
	return &insCollect
}
