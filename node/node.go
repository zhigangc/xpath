package node

import "htmlparser"

type Node struct {
	node *html.Node
	attribute *html.Attribute
}

func (n *Node) Node() *html.Node {
	if n.attribute != nil {
		return nil
	}
	return n.node
}

func (n *Node) Attribute() *html.Attribute {
	return n.attribute
}

func (n *Node) Parent() *html.Node {
	if n.attribute != nil {
		return n.node
	}
	return n.node.Parent
}

func (n *Node) String() string {
	if n.attribute != nil {
		return n.attribute.Val
	}

	nn := n.Node()
	if nn == nil {
		return ""
	}
	return StrVal(nn)
}

func StrVal(n *html.Node) string {
	return ""
}