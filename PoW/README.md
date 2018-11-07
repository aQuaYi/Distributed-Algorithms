# PoW (Proof of Work) 工作量证明共识算法

本目录下的代码由 [Jeiwan](https://github.com/Jeiwan) 的 [blockchain_go](https://github.com/Jeiwan/blockchain_go) 重构而来。

TODO: 本目录下的代码着重于演示 Pow 算法，会弱化其余部分。

## 哈希计算

## math/big 库

## sha256 库

Go 的标准库 [sha256](https://golang.org/pkg/crypto/sha256/) 实现了[安全散列算法-2（SHA-2）](https://zh.wikipedia.org/zh-cn/SHA-2) 中的 SHA-224 和 SHA-256 算法标准。

`sha256.Sum256(data)` 会返回一个 [32]byte 数组作为 data 的校验和。

关于校验和有两个基本点：

- data 的微小改变会带来校验和的极大变化
- 校验和无法反求 data

## encoding/binary

## Merkle Tree