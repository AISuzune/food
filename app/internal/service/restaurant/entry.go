package restaurant

type Group struct{}

// insInfo 创建一个餐厅信息的实例
var insInfo = SInfo{}

func (g *Group) Info() *SInfo {
	return &insInfo
}
