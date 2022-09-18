package rbTree

import (
	"container/list"
	"fmt"
)

/*
	rb树节点
*/
//左旋
func (ins *rb_node) rb_rotate_left() {
	right := ins.rb_right

	ins.rb_right = right.rb_left
	if nil != ins.rb_right {
		ins.rb_right.rb_parent = ins
	}
	right.rb_left = ins

	right.rb_parent = ins.rb_parent
	if nil != right.rb_parent {
		if ins == ins.rb_parent.rb_left {
			ins.rb_parent.rb_left = right
		} else {
			ins.rb_parent.rb_right = right
		}
	} else {
		//ins是根节点
		ins.root.node = right
	}
	ins.rb_parent = right
}

//右旋
func (ins *rb_node) rb_rotate_right() {
	left := ins.rb_left

	ins.rb_left = left.rb_right
	if nil != ins.rb_left {
		left.rb_right.rb_parent = ins
	}
	left.rb_right = ins

	left.rb_parent = ins.rb_parent
	if nil != left.rb_parent {
		if ins == ins.rb_parent.rb_left {
			ins.rb_parent.rb_left = left
		} else {
			ins.rb_parent.rb_right = left
		}
	} else {
		ins.root.node = left
	}
	ins.rb_parent = left
}

//后继
func (ins *rb_node) rb_next() *rb_node {
	if nil != ins.rb_right {
		//节点有右子树
		node := ins.rb_right
		for nil != node.rb_left {
			node = node.rb_left
		}
		return node
	}
	/*
		节点无右子树：
			若其=父节点左节点则父节点就是下一个；
			若其=父节点右节点则向上回溯直至到根节点或=祖先节点的左节点
	*/
	for nil != ins.rb_parent && ins == ins.rb_parent.rb_right {
		ins = ins.rb_parent
	}
	return ins.rb_parent
}

//前驱
func (ins *rb_node) rb_prev() *rb_node {
	if nil != ins {
		//节点有左子树
		for nil != ins.rb_left {
			ins = ins.rb_right
		}
		return ins
	}
	/*
		节点无左子树：
			若其=父节点右节点则父节点就是下一个；
			若其=父节点左节点则向上回溯直至到根节点或=祖先节点的右节点
	*/
	for nil != ins.rb_parent && ins == ins.rb_parent.rb_left {
		ins = ins.rb_parent
	}
	return ins.rb_parent
}

//关联父节点
func (ins *rb_node) rb_link(obj RbNodeItem, parent *rb_node, root *Rb_tree) {
	ins.rb_parent = parent //此时直接就挂上
	ins.rb_color = RB_RED  //初始=r（尽量不违反RB树约束）
	ins.rb_left = nil
	ins.rb_right = nil
	ins.data = obj
	ins.root = root
}

//插入后的树调整
func (ins *rb_node) rb_insert() {
	var (
		parent           = ins.rb_parent
		gparent *rb_node = nil
	)
	//fmt.Printf("%v %d\n", ins.root.node.data.String(), ins.rb_color)
	//插入节点=红色，故父节点若=黑色则对平衡无影响
	for nil != parent && RB_RED == parent.rb_color { //父节点=红色,可推断肯定有祖父节点（不为空&红色）
		gparent = parent.rb_parent

		if parent == gparent.rb_left {
			//父节点是左儿子
			{
				//若有叔叔节点，首次进入理论上应该=红色，因为此时父节点下没有黑节点（后续就不保证了）
				uncle := gparent.rb_right
				if nil != uncle && uncle.rb_color == RB_RED { //有叔叔节点&红色->出现双红(ins和父层节点)
					//将父节点和叔叔节点与祖父节点的颜色互换-->消除双红
					uncle.rb_color = RB_BLACK
					parent.rb_color = RB_BLACK
					gparent.rb_color = RB_RED
					//向上回溯（消除后可能导致上层节点出现双红）
					ins = gparent
					/*
						这里会疑惑为啥不是parent=gparent(而是直接再跳一层)
						原因是若gparent的父节点=红色，根据RB树约束gparent的兄弟节点肯定=黑色（即该层无需调整）
					*/
					parent = ins.rb_parent
					continue
				}
			}
			if parent.rb_right == ins {
				//祖父节点/父节点/新节点不处于一条斜线上，父节点左旋，调整到一条斜线
				parent.rb_rotate_left()
				tmp := parent
				parent = ins
				ins = tmp
			}
			//祖父节点/父节点/新节点处于一条斜线上，祖父节点和父节点互换颜色，祖父节点右旋
			parent.rb_color = RB_BLACK
			gparent.rb_color = RB_RED
			gparent.rb_rotate_right()
		} else {
			//父节点是右儿子（上面分支的镜像）
			{
				uncle := gparent.rb_left
				if nil != uncle && uncle.rb_color == RB_RED {
					uncle.rb_color = RB_BLACK
					parent.rb_color = RB_BLACK
					gparent.rb_color = RB_RED
					ins = gparent
					parent = ins.rb_parent
					continue
				}
			}
			if parent.rb_left == ins {
				parent.rb_rotate_right()
				tmp := parent
				parent = ins
				ins = tmp
			}
			parent.rb_color = RB_BLACK
			gparent.rb_color = RB_RED
			gparent.rb_rotate_left()
		}
	}
	ins.root.node.rb_color = RB_BLACK //根节点=黑色 always
}

//已经剔除后的树调整
func _rb_erase_adjust(node, parent *rb_node) {
	/*
		对于场景2（node是一个红色叶子节点）而言，直接将颜色=黑色即可
	*/
	var other *rb_node //涉及变化的节点的
	for (nil == node || node.rb_color == RB_BLACK) && node != node.root.node {
		/*
			当节点=根节点时直接跳出并=黑色
			对于场景1涉及场景：
				(1).other=黑色（可能有红子节点），parent=黑色----->包含2层黑色
					other=红色（分支的黑色节点数减1），导致不平衡，触发向上回溯
						parent指向其父节点，触发场景1（黑色数+1）
				(2).other=黑色（可能有红子节点），parent=红色----->包含1层黑色
					无子节点：
						other和parent互换颜色即可
					有子节点：
						1个：
							<1>parent/other和子节点是一条线：
								other使用parent的颜色；parent和子节点=黑色；parent左旋
							<2>parent/other和子节点非一条线
								other和子节点互换颜色，other右旋--->成为一条线((2).<1>)
						2个：
							按1个中的<1>执行
				(3).other=红色，可知有2个黑色子节点（父节点=黑色）----->包含2层黑色
					other和parent互换颜色;parent左旋-->一个分支转到1.(2)场景
		*/
		if parent.rb_left == node { //父节点左分支=空（首次进入）---------->父节点右倾
			other = parent.rb_right
			if other.rb_color == RB_RED { //右节点=红色，可通过父节点左旋达到本层平衡
				//1.(3)场景
				other.rb_color = RB_BLACK
				parent.rb_color = RB_RED
				parent.rb_rotate_left()
				other = parent.rb_right
				//调整（other=黑色;父节点=红色;父节点左旋）后，左分支转为场景=1.(2)		缩小范围
			}
			if (nil == other.rb_left || other.rb_left.rb_color == RB_BLACK) &&
				(nil == other.rb_right || other.rb_right.rb_color == RB_BLACK) {
				/*
					场景1.(1)下other节点无红子节点-->other=红色,使本分支保持平衡，但导致上层不平衡，要向上回溯（父节点（黑色）不会退出循环）
						扩大范围
					场景1.(2)下other节点无红子节点-->other和父节点互换颜色（实现：other=红色;父节点（红色）退出循环后在兜底处=黑色）
					回溯过程中存在父节点左旋（上面的分支）后的上层平衡而将另一分支的子节点转移过来
						若该子节点中的左右节点=黑色，可尝试将该子节点=红色（黑色数减1），而使子分支平衡
				*/
				other.rb_color = RB_RED
				node = parent
				parent = node.rb_parent
			} else { //该分支说明父节点右倾
				/*
					场景1.(2)下other节点有红子节点
				*/
				if nil == other.rb_right || other.rb_right.rb_color == RB_BLACK { //该分支说明other左倾--->调到一条线（与整体倾向一致）
					/*
						只有左子节点（parent=红,other=黑,左子=红），调成为一条线
					*/
					if nil != other.rb_left {
						//other节点若有左子节点,则=黑色
						other.rb_left.rb_color = RB_BLACK
					}
					other.rb_color = RB_RED
					other.rb_rotate_right()
					other = parent.rb_right
					//调整后-->场景1.(2)且other无子节点
				}
				other.rb_color = parent.rb_color
				parent.rb_color = RB_BLACK
				if nil != other.rb_right {
					other.rb_right.rb_color = RB_BLACK
				}
				parent.rb_rotate_left()
				node = node.root.node //触发退出循环
				break
			}
		} else {
			//上面分支的镜像
			other = parent.rb_left
			if other.rb_color == RB_RED {
				other.rb_color = RB_BLACK
				parent.rb_color = RB_RED
				parent.rb_rotate_right()
				other = parent.rb_left
			}
			if (nil == other.rb_left || other.rb_left.rb_color == RB_BLACK) &&
				(nil == other.rb_right || other.rb_right.rb_color == RB_BLACK) {
				other.rb_color = RB_RED
				node = parent
				parent = node.rb_parent
			} else {
				if nil == other.rb_left || other.rb_left.rb_color == RB_BLACK {
					if nil != other.rb_right {
						other.rb_right.rb_color = RB_BLACK
					}
					other.rb_color = RB_RED
					other.rb_rotate_left()
					other = parent.rb_left
				}
				other.rb_color = parent.rb_color
				parent.rb_color = RB_BLACK
				if nil != other.rb_left {
					other.rb_left.rb_color = RB_BLACK
				}
				parent.rb_rotate_right()
				node = node.root.node //触发退出循环
				break
			}
		}
	}
	if nil != node {
		node.rb_color = RB_BLACK
	}
}

func (ins *rb_node) rb_erase() {
	var (
		child  *rb_node = nil
		parent *rb_node = nil
		color  int
	)

	if nil == ins.rb_left {
		//无左节点
		child = ins.rb_right //此时可能有右节点（左无右有）or无子节点
	} else if nil == ins.rb_right {
		//无右节点
		child = ins.rb_left //此时有左节点（左有右无）
	} else {
		//左右节点都有，则用后继节点替换删除节点
		old := ins
		//找到后继节点
		ins = ins.rb_right
		for nil != ins.rb_left {
			ins = ins.rb_left
		}
		child = ins.rb_right
		parent = ins.rb_parent
		color = ins.rb_color

		//把后继节点摘出
		if nil != child {
			//存在右子树
			child.rb_parent = parent
		}
		if nil != parent {
			if ins == parent.rb_left {
				parent.rb_left = child
			} else {
				parent.rb_right = child
			}
		} else {
			//根节点
			ins.root.node = child
		}

		if ins.rb_parent == old {
			//后继节点的父节点正好是要删除的节点时，需要进行调整：替换后后继节点将成为父节点
			parent = ins
		}
		//用后继节点替换要删除的节点，把要删除的节点摘出
		ins.rb_parent = old.rb_parent
		ins.rb_color = old.rb_color //包括颜色，因此平衡判断时考虑的是后继节点的影响
		ins.rb_right = old.rb_right
		ins.rb_left = old.rb_left
		if nil != old.rb_parent {
			if old == old.rb_parent.rb_left {
				old.rb_parent.rb_left = ins
			} else {
				old.rb_parent.rb_right = ins
			}
		} else {
			//根节点
			ins.root.node = ins
		}

		old.rb_left.rb_parent = ins
		if nil != old.rb_right {
			old.rb_right.rb_parent = ins
		}
		goto COLOR
	}

	parent = ins.rb_parent
	color = ins.rb_color

	//无子节点or只有一个子节点场景下，剔除要删除的节点
	if nil != child {
		child.rb_parent = parent
	}
	if nil != parent {
		if ins == parent.rb_left {
			parent.rb_left = child
		} else {
			parent.rb_right = child
		}
	} else {
		//根节点
		ins.root.node = child
	}

COLOR:
	//若有且只有1个子节点时，根据RB树定义被删除的节点=黑色
	if RB_BLACK == color {
		/*
			要删除的节点=黑色，会破坏RB树特性，需调整
			涉及场景：
				场景1：删除的节点无子节点--->
					父节点（r or b）且另一分支肯定包含有黑节点
					子节点=空
				场景2：删除的节点有一个子节点（=红色）--->
					父节点（r or b）且另一分支肯定包含有黑节点
					子节点（红色）
				删除的节点有2个子节点，在后继节点替换后--->
					上述2种情况都存在
		*/
		_rb_erase_adjust(child, parent)
	}
}

/*
	rb树
*/
//rb树第一个节点（min）
func (ins *Rb_tree) rb_first() *rb_node {
	node := ins.node
	if nil == node {
		return nil
	}
	for nil != node.rb_left {
		node = node.rb_left
	}
	return node
}

//rb树最后一个节点（max）
func (ins *Rb_tree) rb_last() *rb_node {
	node := ins.node
	if nil == node {
		return nil
	}
	for nil != node.rb_right {
		node = node.rb_right
	}
	return node
}

func (ins *Rb_tree) search(obj RbNodeItem) (*rb_node, int, bool) {
	var (
		loc             = ROOT
		cur             = ins.node
		parent *rb_node = nil
	)
	for nil != cur {
		result := cur.data.Compare(obj)
		if 0 == result {
			return parent, loc, true
		} else {
			parent = cur
			if 0 > result {
				loc = RIGHT_CHILD
				cur = cur.rb_right
			} else {
				loc = LEFT_CHILD
				cur = cur.rb_left
			}
		}
	}
	return parent, loc, false
}

func (ins *Rb_tree) Insert(obj RbNodeItem) bool {
	parent, loc, exist := ins.search(obj)
	if exist {
		return false
	}
	if nil == parent && ROOT != loc {
		panic("异常，不应该出现")
	}
	node := new(rb_node)
	node.rb_link(obj, parent, ins)
	if nil == parent && ROOT == loc {
		//根节点
		node.rb_color = RB_BLACK
		ins.node = node
		return true
	}
	if LEFT_CHILD == loc {
		parent.rb_left = node
	} else {
		parent.rb_right = node
	}
	//fmt.Printf("%v %d\n", ins.node.data.String(), ins.node.rb_color)
	node.rb_insert()
	return true
}

func (ins *Rb_tree) Delete(obj RbNodeItem) {
	parent, loc, exist := ins.search(obj)
	if !exist {
		return
	}

	var node *rb_node
	if LEFT_CHILD == loc {
		node = parent.rb_left
	} else {
		node = parent.rb_right
	}
	if nil != node {
		node.rb_erase()
	}
}

func (ins *Rb_tree) PreTraverse() []string {
	out := make([]string, 0)
	first := ins.rb_first()
	for nil != first {
		//fmt.Printf("%+v \n", first)
		out = append(out, fmt.Sprintf("%s with %v", first.data.String(), first.rb_color))
		first = first.rb_next()
	}
	return out
}

type rb_node_with_level struct {
	*rb_node
	level uint32
}

func (ins *Rb_tree) LevelTraverse() [][]string {
	if nil == ins.node {
		//空
		return nil
	}
	var (
		maxLevel = uint32(0)
		queue    = list.New()
		mNodes   = make(map[uint32][]string)
	)
	queue.PushBack(&rb_node_with_level{
		rb_node: ins.node,
		level:   1,
	})
	for e := queue.Front(); nil != e; e = e.Next() {
		cur := e.Value.(*rb_node_with_level)
		if nil != cur.rb_left {
			queue.PushBack(&rb_node_with_level{
				rb_node: cur.rb_left,
				level:   cur.level + 1,
			})
		}
		if nil != cur.rb_right {
			queue.PushBack(&rb_node_with_level{
				rb_node: cur.rb_right,
				level:   cur.level + 1,
			})
		}
		vNode, prs := mNodes[cur.level]
		if !prs {
			vNode = make([]string, 0, 1)
		}
		vNode = append(vNode, fmt.Sprintf("%s with %v", cur.data.String(), cur.rb_color))
		mNodes[cur.level] = vNode
		if cur.level > maxLevel {
			maxLevel = cur.level
		}
	}

	out := make([][]string, 0, maxLevel)
	for idx := uint32(1); idx <= maxLevel; idx++ {
		vNode, prs := mNodes[idx]
		if !prs {
			out = append(out, []string{})
			continue
		}
		out = append(out, vNode)
	}
	return out
}
