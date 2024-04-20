package recipe

type Group struct{}

// insInfo 创建一个菜谱信息的实例
var insInfo = SInfo{}

func (g *Group) Info() *SInfo {
	return &insInfo
}
