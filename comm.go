package rbtree

const (
	RB_RED   = 0
	RB_BLACK = 1
	//
	ROOT        = 0
	LEFT_CHILD  = 1
	RIGHT_CHILD = 2
)

type RbNodeItem interface {
	Compare(RbNodeItem) int //=0表示匹配；>1表示大于（找左子树）；<1表示小于（找右子树）
	String() string
}

type rb_node struct {
	rb_parent *rb_node
	rb_right  *rb_node
	rb_left   *rb_node
	rb_color  int //r or b
	root      *Rb_tree
	data      RbNodeItem
}

type Rb_tree struct {
	node *rb_node
}
