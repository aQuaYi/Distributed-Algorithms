# Mutual Exclusion Algorithm Demo

使用 Go 语言实现了 Lamport 在论文 《Time, Clocks and the Ordering of Events in a Distributed System》中提到的 Mutual Exclusion 算法

## 问题

system 由多个 process 组成，却只有一个 resource 可供使用。resource 最多只能被一个 process 占用。由于所有的 process 都是平等的，因此，程序中没有负责调度占用 resource 的功能。需要靠算法满足以下功能：

1. 占用 resource 的 process，在别的 process 占用 resource 前，一定要释放资源。
1. process 占用 resource 的顺序，要和他们申请占用 resource 的顺序必须一致。
1. 如果 process 一定会释放 resource，那么，所有占用 resource 的申请，一定会被满足。

为了简化问题，还存在以下假设：

1. 任意两个 process 都可以直接相互发送消息
1. 对于任意两个 process Pi 和 Pj 而言，从 Pi 发往 Pj 的消息，满足先发送先到达的原则
1. process 间发送的消息，一定会收到

## 思路

## 代码说明
