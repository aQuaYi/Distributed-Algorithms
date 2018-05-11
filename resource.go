package mutual

import (
	"fmt"
)

const (
	// NULL 表示没有赋予任何人
	NULL = -1
)

var (
	rsc         *resource // 全局变量，随时随地都可以访问
	occupyOrder []int     // rsc 被占用的顺序
)

func init() {
	rsc = &resource{
		grantedTo: NULL,
	}
}

type resource struct {
	grantedTo int
}

func (r *resource) occupy(p int) {
	if r.grantedTo != NULL {
		msg := fmt.Sprintf("资源正在被 P%d 占据，P%d 却想获取资源。", r.grantedTo, p)
		panic(msg)
	}
	r.grantedTo = p
	occupyOrder = append(occupyOrder, p)
}

func (r *resource) release(p int) {
	if r.grantedTo != p {
		msg := fmt.Sprintf("P%d 想要释放正在被 P%d 占据的资源。", p, r.grantedTo)
		panic(msg)
	}
	r.grantedTo = NULL
}
