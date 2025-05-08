package bitcaskgo

// 用户上传的配置文件

type Options struct {
	DirPath string // 数据文件存储的目录路径
	DataFileSize int64 // 数据文件的最大大小，单位字节
	SyncWrites bool // 是否每次写入都进行同步操作（持久化）
}