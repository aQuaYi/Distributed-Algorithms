# 思考问题

1.为什么会出现多种全局排序？请举例说明。

```text
由于不同 process 的 event 可能会有相同的 timestamp，例如 E1 = `<3:7>`, E2 = `<3:5>`。
如果 P7 < P5 的话， E1 => E2。
如果 P5 < P7 的话， E2 => E1。
```

2.真实时间上先 request 的 process 会不会后得到 resource？如果会的话，能不能说明 mutual exclusion 算法失败了？请说明理由。

```text
会，但不是全局排序的失败。

我在编程的时候，就遇到了这个问题。开始以为是程序的 bug 。后来重新阅读的论文，才发现是自己的理解的不够。

首先区分一下时间(time)和时刻(timestamp)，时间是一个物理量，时刻是这个物理量的值，2018年05月15日15:20:55 是现在的时刻。就好像温度是一个物理量，33℃是温度的一个值。但是如果33℃的物体比44℃的物体摸起来要热，只能说明这个物体不是使用同一个温度计测量的温度，并且两个温度计的基准差别还蛮大。
第二，时刻(timestamp)的作用是给 event 一个标记，多个 event 可以利用这个时间标记进行排序。例如，同一天中，E1(15:31:31) 排在 E2(15:31:51) 前面。但这包含了一个隐含前提，这两个 event 的时刻，是由同一个可靠的 clock 标记的。
第三，mutual exclusion 是一个分布式算法。每个 process 都有自己单独的 clock。不同 process 中的 event 的时间标记都是不同的 clock 标记的。考虑到程序运行的速度，这些 clock 与真实时间之间的偏差，绝对不能忽略不计。
第四，为了 process 间的局部排序，引入了 message 机制，并制定了 lamport timestamp 规则。为了全局排序，再引入 process 排序。

再解释一下题意，存在一个观察者，拿着同一个 clock 去分别标记每个 process 的 request，结果发现某个先标记的 request 却后得到了 resource。
这不能说明 mutual exclusion 算法失败的原因是，这个算法就是为了解决分布式系统中，不可能存在同一个 clock 去分别标记每个 process 的 request 的问题而提出的。

如果能像 Google 在 [Spanner](spanner-osdi2012.pdf) 里面，引入 True Time 一样，使得各个 process 的 clock 之间的偏差，相对于程序的速度可以忽略不计。就可以保证真实时间上先 request 的 process 先占用 resource。那样的话，也不需要 mutual exclusion 算法了。
```
