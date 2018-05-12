package mutual

func (p *process) append(r *request) {
	p.requestQueue = append(p.requestQueue, r)
	debugPrintf("[%d]P%d append %s, MRT=%d, RT%v, RQ%v", p.clock.getTime(), p.me, r, p.minReceiveTime, p.receiveTime, p.requestQueue)
}

func (p *process) delete(r *request) {
	i := 0
	for p.requestQueue[i] != r {
		i++
	}
	last := len(p.requestQueue) - 1

	// 删除的时候，需要保持 requestQueue 的顺序
	copy(p.requestQueue[i:], p.requestQueue[i+1:])

	p.requestQueue = p.requestQueue[:last]

	debugPrintf("[%d]P%d delete %s, MRT=%d, RT%v, RQ%v", p.clock.getTime(), p.me, r, p.minReceiveTime, p.receiveTime, p.requestQueue)

	// p.requestQueue 变化时，都需要检查是否符合了 rule5
	p.toCheckRule5Chan <- struct{}{}
}
