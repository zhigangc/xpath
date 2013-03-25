package xnode

import "htmlparser"

type XPathNode struct {
	node *html.Node
	attribute *html.Attribute
}

func (n *XPathNode) Node() *html.Node {
	if n.attribute != nil {
		return nil
	}
	return n.node
}

func (n *XPathNode) Attribute() *html.Attribute {
	return n.attribute
}

func (n *XPathNode) Parent() *html.Node {
	if n.attribute != nil {
		return n.node
	}
	return n.node.Parent
}

func (n *XPathNode) String() string {
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