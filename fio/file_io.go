package fio

import "os"

// FileIO 标准文件系统 IO
type FileIO struct {
	fd *os.File // 系统文件描述符

}

// NewFileIOManager 初始化标准文件 FileIO 实例
func NewFileIOManager(fileName string) (*FileIO, error) {
	fd, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_RDWR|os.O_APPEND, // 打开文件，如果不存在则创建，如果存在则打开，并且赋予读取权限
		DataFilePerm,                      // 文件权限，0644 表示文件所有者有读写权限，其他用户只有读权限
	)
	if err != nil {
		return nil, err
	}
	return &FileIO{fd: fd}, nil
}

func (fio *FileIO) Read(b []byte, offset int64) (int, error) {
	return fio.fd.ReadAt(b, offset)
}

func (fio *FileIO) Write(b []byte) (int, error) {
	return fio.fd.Write(b)
}

func (fio *FileIO) Sync() error {
	return fio.fd.Sync()
}

func (fio *FileIO) Close() error {
	return fio.fd.Close()
}
