package collx

// INode 其他的结构体想要生成菜单树，直接实现这个接口
type INode interface {
	// GetId获取id
	GetId() int
	// GetPid 获取父id
	GetPid() int
	// IsRoot 判断当前节点是否是顶层根节点
	IsRoot() bool

	SetChildren(childern any)
}

// INodes 实现 sort.Interface 的 INode 切片。
type INodes []INode

func (nodes INodes) Len() int {
	return len(nodes)
}
func (nodes INodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}
func (nodes INodes) Less(i, j int) bool {
	return nodes[i].GetId() < nodes[j].GetId()
}

// GenerateTree 自定义的结构体实现 INode 接口后调用此方法生成树结构
// nodes 需要生成树的节点
// 返回生成后的树结构对象
func GenerateTree(nodes []INode) (trees []INode) {
	trees = []INode{}
	// 定义顶层根和子节点
	var roots, childs []INode
	for _, v := range nodes {
		if v.IsRoot() {
			// 判断顶层根节点
			roots = append(roots, v)
		}
		childs = append(childs, v)
	}

	for _, v := range roots {
		// 递归
		setChildren(v, childs)
		trees = append(trees, v)
	}
	return
}

// setChildren 递归地把 nodes 中 parent 的子节点挂到 parent 的 Children 上
// parent 当前父节点
// nodes 待挂载的全部节点
func setChildren(parent INode, nodes []INode) {
	children := []INode{}
	for _, v := range nodes {
		if v.IsRoot() {
			// 如果当前节点是顶层根节点就跳过
			continue
		}
		if parent.GetId() == v.GetPid() {
			children = append(children, v)
		}
	}
	if len(children) == 0 {
		return
	}

	parent.SetChildren(children)
	for _, c := range children {
		setChildren(c, nodes)
	}
}
