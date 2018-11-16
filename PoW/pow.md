# PoW

## 为什么需要 PoW

区块链是典型的[分布式系统](https://wiki.mbalib.com/wiki/%E5%88%86%E5%B8%83%E5%BC%8F%E7%B3%BB%E7%BB%9F)。每个节点都想要获取区块的写入权，以便获得挖矿奖励。为了公平地分配写入权，运用 PoW 机制，第一个完成指定工作的节点，会得到写入权。

## 什么是工作量证明

1. 把准备写入的区块对象的关键属性，转换成 []byte 格式的 bytes
2. 在 bytes 尾部添加计数器 nonce
3. 获取 bytes 的 hash 值
4. 把 hash 值按照大端法解释成整数 hashInt ,
  4.1. 如果， hashInt < 事先指定的 target，则完成工作量
  4.2. 否则， nonce++，回到步骤 2

优点

1. 找到 nonce 很难，验证 nonce 却很容易。

缺点

1. 找到 nonce 很耗时，耗电量太大，对环境不友好。

## 细节

// TODO: 补充完整细节

```golang
func main() {
    a := big.NewInt(0x1234)
    fmt.Printf("0x%x, %v\n", a, a)
    // 设置 a 为 0x5678
    a.SetBytes([]byte{0x56, 0x78})
    fmt.Printf("0x%x, %v\n", a, a)
}
```

[点击运行](https://play.golang.org/p/EHCDZT1zadc)
