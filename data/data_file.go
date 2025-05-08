package data

import "bitcask-go/fio"

// 数据文件
type DataFile struct {
	FileId    uint32
	WritePos  int64         // 写入位置（偏移量）
	IoManager fio.IOManager // 文件IO管理器
}

// 打开数据文件，返回数据文件对象
func OpenDataFile(dirPath string, fileId uint32) (*DataFile, error) {
	return nil, nil
}

func (df *DataFile) WriteOffset(buf []byte, offset int64) error {
	return nil
}

func (df *DataFile) Sync() error {
	return nil
}
