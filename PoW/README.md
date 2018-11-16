# PoW (Proof of Work) 工作量证明共识算法

本目录下的代码由 [Jeiwan](https://github.com/Jeiwan) 的 [blockchain_go](https://github.com/Jeiwan/blockchain_go) 重构而来。

TODO: 本目录下的代码着重于演示 Pow 算法，会弱化其余部分。

## 散列函数 (Hash function)

[散列函数](https://zh.wikipedia.org/zh-cn/%E6%95%A3%E5%88%97%E5%87%BD%E6%95%B8)具有以下特点:

- 正向快速：可以很快地计算出 hash 值。
- 逆向困难：给定 hash 值，基本不可能逆推出明文。
- 输入敏感：输入的细微改变，输出的 hash 值也会明显不同。
- 冲突避免：几乎不可能找到两个输入值，具有相同的输出。

比特币中使用的 sha256 散列算法，意味着其输出的 hash 值的长度为 256 bit。 // TODO: 完善此处的说明

## PoW

// TODO: 完整描述 PoW 的工作机制

## 哈希计算

### math/big 库

### sha256 库

Go 的标准库 [sha256](https://golang.org/pkg/crypto/sha256/) 实现了[安全散列算法-2（SHA-2）](https://zh.wikipedia.org/zh-cn/SHA-2) 中的 SHA-224 和 SHA-256 算法标准。

`sha256.Sum256(data)` 会返回一个 [32]byte 数组作为 data 的校验和。

关于校验和有两个基本点：

- data 的微小改变会带来校验和的极大变化
- 校验和无法反求 data

<!-- ### encoding/binary -->

## UTXO (Unspent Transaction Output)

## Merkle Tree

## 公钥（public key）和签名（signature）

## 参考链接

- <https://blog.csdn.net/asdzheng/article/details/70226007>