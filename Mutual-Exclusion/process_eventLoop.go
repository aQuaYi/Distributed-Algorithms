package mutual

func eventLoop(p *process) {
	debugPrintf("[%d]P%d 启动 eventLoop", p.clock.getTime(), p.me)

	go func() {

		for {
			p.clock.tick()
			select {
			case msg := <-p.chans[p.me]:
				p.handleMsg(msg)
			case <-p.requestChan:
				p.handleRequest()
			case <-p.toCheckRule5Chan:
				p.handleCheckRule5()
			}
		}

	}()
}

func (p *process) handleMsg(msg *message) {
	debugPrintf("[%d]P%d receive %s RQ%s", p.clock.getTime(), p.me, msg, p.requestQueue)

	// 接收到了一个新的消息
	// 根据 IR2
	// process 的 clock 需要根据 msg.time 进行更新
	// 无论 msg 是什么类型的消息
	nowTime := p.clock.update(msg.timestamp)

	p.receiveTime[msg.from] = nowTime

	p.updateMinReceiveTime()

	r := msg.request

	switch msg.msgType {
	case requestResource:
		// 根据 Rule2
		p.push(r)
	case releaseResource:
		// 根据 Rule4
		p.pop(r)
	}

	// NOTICE: 与论文中不同
	// 我总是发送 acknowledgement 信息
	if msg.msgType != acknowledgment {
		p.chans[msg.from] <- newMessage(acknowledgment, nowTime, p.me, nil)
	}

	// 每次收到了消息，都会触发检查，是否已经满足 Rule5
	go func() {
		p.toCheckRule5Chan <- struct{}{}
	}()

	debugPrintf("[%d]P%d finish handleMsg %s RQ%s", p.clock.getTime(), p.me, msg, p.requestQueue)
}

func (p *process) handleRequest() {
	timestamp := p.clock.getTime()
	r := &request{
		timestamp: timestamp,
		process:   p.me,
	}

	debugPrintf("[%d]P%d 开始 handleRequest，request message%s", timestamp, p.me, r)

	// 根据 Rule1
	//
	// Rule1.0 给其他的 process 发消息
	for i := range p.chans {
		if i == p.me {
			continue
		}
		p.chans[i] <- newMessage(requestResource, timestamp, p.me, r)
	}

	// Rule1.1 把 r 放入自身的 request queue
	p.push(r)
}

func (p *process) handleRelease() {
	timestamp := p.clock.getTime()
	debugPrintf("[%d]P%d 开始 handleRelease, request queue %v", timestamp, p.me, p.requestQueue)

	req := p.requestQueue[0]

	// 根据 Rule 3
	//
	// Rule3.0 释放资源
	p.resource.release(req)
	// 标记自己已释放
	p.isOccupying = false

	// Rule3.1 在 requestQueue 中删除 req
	p.pop(req)

	// Rule3.2 给其他的 process 发消息
	for i := range p.chans {
		if i == p.me {
			continue
		}
		p.chans[i] <- newMessage(releaseResource, timestamp, p.me, req)
	}
}

func (p *process) handleCheckRule5() {
	debugPrintf("[%d]P%d to check Rule5", p.clock.getTime(), p.me)

	if len(p.requestQueue) > 0 && // p.requestQueue 中还有元素
		p.requestQueue[0].process == p.me && // 排在首位的 repuest 是 p 自己的
		p.requestQueue[0].timestamp < p.minReceiveTime && // p 在 request 后，收到过所有其他 p 的回复
		!p.isOccupying { // 不能是正占用的资源

		debugPrintf("[%d]P%d 满足 Rule5 MRT=%d RT%v PQ%v", p.clock.getTime(), p.me, p.minReceiveTime, p.receiveTime, p.requestQueue)
		p.handleOccupy()

	}

	debugPrintf("[%d]P%d 不满足 Rule5 MRT=%d RT%v PQ%v", p.clock.getTime(), p.me, p.minReceiveTime, p.receiveTime, p.requestQueue)
}
func (p *process) handleOccupy() {
	req := p.requestQueue[0]
	debugPrintf("[%d]P%d handleOccupy %s request queue %v", p.clock.getTime(), p.me, req, p.requestQueue)

	p.isOccupying = true

	p.resource.occupy(req)

	// 假设条件，process 不会永远占用 resource
	go func() {
		randSleep()
		// p.releaseChan <- struct{}{}
		p.handleRelease()
	}()
}
