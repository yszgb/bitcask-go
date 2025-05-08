package bitcaskgo

import (
	"bitcask-go/data"
	"bitcask-go/index"
	"sync"
)

// bitcask 存储引擎实例
type DB struct {
	options    Options
	mu         *sync.RWMutex
	activeFile *data.DataFile   // 当前活跃的文件
	olderFiles []*data.DataFile // 旧的文件列表，只读
	index      index.Indexer    // 内存索引，内存中保存的索引数据结构
}

// 写入
func (db *DB) Put(key []byte, value []byte) error {
	// 判断 key 是否有效
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	// 构造 LogRecord 结构体
	log_record := &data.LogRecord{
		Key:   key,
		Value: value,
		Type:  data.LogRecordNormal,
	}

	pos, err := db.appendLogRecord(log_record)
	if err != nil {
		return err
	}

	// 更新索引
	if ok:=db.index.Put(key, pos);!ok{
		 return ErrIndexUpdateFailed
	}

	return nil
}

// 根据 key 获取数据
func (db *DB) Get(key []byte) ([]byte, error) {

}

// 追加写数据到活跃文件
func (db *DB) appendLogRecord(logRecord *data.LogRecord) (*data.LogRecord, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 初始化活跃文件
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	// 写入数据
	encRecord, size := data.EncodeLogRecord(logRecord)
	// 到达阈值，关闭当前文件，打开新的文件
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		// 先持久化
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}

		// 转为旧文件
		db.olderFiles[db.activeFile.FileId] = db.activeFile

		// 打开新的文件
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	// 写入数据
	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}

	// 根据用户配置决定是否持久化
	if db.options.SyncWrites {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	// 更新写入位置，构造索引
	pos := &data.LogRecordPos{Fid: db.activeFile.FileId, Offset: writeOff}
	return pos, nil
}

// 初始化活跃文件
//
// 访问此方法前必须持有互斥锁
func (db *DB) setActiveDataFile() error {
	var initialFileId uint32 = 0
	if db.activeFile != nil {
		initialFileId = db.activeFile.FileId + 1
	}
	dataFile, err := data.OpenDataFile(db.options.DirPath, initialFileId)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}
