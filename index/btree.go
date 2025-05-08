package index

import (
	"bitcask-go/data"
	"github.com/google/btree"
	"sync"
)

// BTree: 一个平衡树的实现，每个节点可以有多个子节点，每个子节点可以有多个子节点，以此类推，直到叶子节点
// 每个节点可以存储多个键值对，每个键值对的键是唯一的，值是任意的
// 每个节点的键值对是按照键的大小排序的，左子节点的键值对的键小于父节点的键值对的键，右子节点的键值对的键大于父节点的键值对的键

// Btree 索引，封装 Google 的 Btree 库
type Btree struct {
	// 查看源码可知，写操作并发不安全，需要加锁
	// 读操作并发安全
	tree *btree.BTree
	lock *sync.RWMutex // 这是一个读写锁，读锁可以并发，写锁不可以并发
}

// 根据 index.go 定义出的 Indexer 接口，实现 Btree 索引的 Put、Get、Delete 方法
func (bt *Btree) Put(key []byte, pos *data.LogRecordPos) bool {
	// bt.tree.ReplaceOrInsert()
	// 需要传入 Item，源码中是一个抽象接口，需要实现 Less 方法
	// index.go 中实现
	it := &Item{key: key, pos: pos}
	bt.lock.Lock()
	bt.tree.ReplaceOrInsert(it)
	bt.lock.Unlock()
	return true
}

// NewBtree 初始化 Btree 索引结构
func NewBtree() *Btree {
	return &Btree{
		tree: btree.New(32), // 控制叶子节点数量
		// 可以考虑让用户进行控制
		lock: new(sync.RWMutex),
	}
}

func (bt *Btree) Get(key []byte) *data.LogRecordPos {
	it := &Item{key: key}
	btreeItem := bt.tree.Get(it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos // 类型断言，将 btreeItem 转换为 *Item 类型，然后返回 pos
	// 为何返回pos？因为在 Item 结构体中，pos 是一个指针，所以需要返回指针
}

func (bt *Btree) Delete(key []byte) bool {
	it := &Item{key: key}
	bt.lock.Lock()
	oldItem := bt.tree.Delete(it)
	bt.lock.Unlock()
	return oldItem != nil
}
