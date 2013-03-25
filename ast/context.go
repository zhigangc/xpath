package ast

import "xpath/node"

type XPathContext struct {
	Node *xnode.XPathNode
	Position int
	Size int
}
 
func NewXPathContext(n *xnode.XPathNode, position, size int) *XPathContext {
	return &XPathContext{Node: n, Position: position, Size:size}
}

const (
	xpath_type_none = iota	  // Unknown type (query failed to compile)
	xpath_type_node_set       // Node set (xpath_node_set)
	xpath_type_number         // Number
	xpath_type_string         // String
	xpath_type_boolean	      // Boolean
)