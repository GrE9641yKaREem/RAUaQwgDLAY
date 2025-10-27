// 代码生成时间: 2025-10-27 13:54:44
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

// AtomicExchange 是一个结构体，用于演示原子交换协议
type AtomicExchange struct {
    // value 存储共享的值
    value int64
    // lock 用于确保线程安全
    lock sync.Mutex
}

// NewAtomicExchange 创建一个新的 AtomicExchange 实例
func NewAtomicExchange(initialValue int64) *AtomicExchange {
    return &AtomicExchange{value: initialValue}
}

// Exchange 原子交换方法，它将当前值替换为 newValue，并返回旧值
func (ae *AtomicExchange) Exchange(newValue int64) (oldValue int64, err error) {
    // 锁定用于保护共享资源
    ae.lock.Lock()
    defer ae.lock.Unlock()
    // 使用原子操作来替换值
    oldValue = atomic.SwapInt64(&ae.value, newValue)
    return oldValue, nil
}

// GetValue 获取当前值（非原子操作，用于演示）
func (ae *AtomicExchange) GetValue() int64 {
    // 锁定用于保护共享资源
    ae.lock.Lock()
    defer ae.lock.Unlock()
    // 返回当前值
    return ae.value
}

func main() {
    // 创建一个新的 AtomicExchange 实例，初始值为 0
    ae := NewAtomicExchange(0)

    // 演示原子交换操作
    oldValue, err := ae.Exchange(10)
    if err != nil {
        fmt.Println("Error during exchange: ", err)
    } else {
        fmt.Printf("Old value: %d, New value: %d
", oldValue, ae.GetValue())
    }
}
