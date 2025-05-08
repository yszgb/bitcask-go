package index

import (
	"bitcask-go/data"
	"bytes"
	"github.com/google/btree"
)

// 索引相关

// Indexer 抽象索引接口。不同数据接口，实现这个接口即可
type Indexer interface {
	// Put 向索引中存储 key 对应的 value，也就是数据的位置信息
	Put(key []byte, pos *data.LogRecordPos) bool // 根据 key 存储 value

	// Get 根据 key 获取到对应的 value
	Get(key []byte) *data.LogRecordPos           // 根据 key 获取到对应的 value

	// Delete 根据 key 删除对应的 value
	Delete(key []byte) bool
}

type Item struct {
	key []byte
	pos *data.LogRecordPos
}

func (ai *Item) Less(bi btree.Item) bool {
	return bytes.Compare(ai.key, bi.(*Item).key) == -1
	// ai是调用者，bi是参数
	// ai.key 小于 bi.key，返回 true，否则返回 false
}
