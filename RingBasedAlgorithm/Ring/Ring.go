package Ring

type LocalNode struct {
	LocalPort string
	NextNode  *LocalNode
}

func NewNode(port string) *LocalNode {
	n := &LocalNode{
		LocalPort: port,
		NextNode:  nil,
	}
	return n
}

func (n *LocalNode) AppendNode(next *LocalNode) {
	n.NextNode = next
}

func CreateRingByPoint(Ports []string) *LocalNode {
	lenP := len(Ports)
	if lenP == 1 {
		return NewNode(Ports[0])
	}
	start := NewNode(Ports[0])
	var end *LocalNode
	i := 1
	for i < lenP {
		n := NewNode(Ports[i])
		AppendNode := start
		for AppendNode.NextNode != nil {
			AppendNode = AppendNode.NextNode
		}
		AppendNode.AppendNode(n)
		if i == lenP-1 {
			end = n
		}
		i++
	}
	end.AppendNode(start)
	return start
}
