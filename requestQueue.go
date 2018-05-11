package mutual

func (p *process) append(r *request) {
	p.requestQueue = append(p.requestQueue, r)
	debugPrintf("[%d]P%d append %s, request queue%v", p.clock.getTime(), p.me, r, p.requestQueue)
	p.clock.tick()
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

	p.clock.tick()

	// p.requestQueue 变化时，都需要检查是否符合了 rule5
	p.toCheckRule5Chan <- struct{}{}
}
