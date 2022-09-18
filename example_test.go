package rbTree_test

import (
	"fmt"
	"strconv"

	"github.com/chris-wangkk/rbTree"
)

type demo struct {
	data int64
}

func (ins *demo) Compare(obj rbTree.RbNodeItem) int {
	demoObj := obj.(*demo)
	if ins.data == demoObj.data {
		return 0
	} else if ins.data > demoObj.data {
		return 1
	} else {
		return -1
	}
}

func (ins *demo) String() string {
	return strconv.FormatInt(ins.data, 10)
}

func Example() {
	ins := new(rbTree.Rb_tree)
	for idx := int64(1); idx <= 100; idx++ {
		ins.Insert(&demo{data: idx * 10})
	}
	fmt.Printf("out:%+v \n", ins.PreTraverse())
	for idx, vNode := range ins.LevelTraverse() {
		fmt.Printf("index:%d		[%+v] \n", idx, vNode)
	}
}
