package fio

import()

const DataFilePerm = 0644 // 数据文件权限，可读可写

// IOManager 抽象 IO 管理接口，可以接入不同的 IO 类型。目前支持标准文件 IO
type IOManager interface{
	// Read 从文件的给定位置读取对应的数据
	Read([]byte,int64)(int,error)

	// Write 写入字节数组到文件中
	Write([]byte)(int,error)

	// Sync 持久化数据到磁盘中
	Sync() error

	// Close 关闭文件
	Close() error
}