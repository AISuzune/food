package user

type Group struct{}

// insUser 是 SUser 的一个实例
var insUser = SUser{}

func (g *Group) User() *SUser {
	return &insUser
}

// insCollect 创建一个收藏的实例
var insCollect = SCollect{}

func (g *Group) Collect() *SCollect {
	return &insCollect
}
