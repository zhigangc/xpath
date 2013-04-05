package ast

import "xpath/node"

type Context struct {
	Node *node.Node
	Position int
	Size int
}
 
func NewContext(n *node.Node, position, size int) *Context {
	return &Context{Node: n, Position: position, Size:size}
}