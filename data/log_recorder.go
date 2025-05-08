package data

type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota // 正常的日志记录类型，用于表示正常的数据写入操作。
	LogRecordDeleted                     // 表示数据被删除的日志记录类型。
)

// 索引内存的数据结构

// 写入到数据文件的记录
// 
// 追加写入，类似日志的格式
type LogRecord struct{
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecordPos 数据内存索引，描述数据在磁盘上的位置
// Pos 指示数据在磁盘上的位置
type LogRecordPos struct {
	Fid    uint32 // 文件id，表示将数据存储到了哪个文件当中
	Offset int64  // 相对于文件的偏移量，代表数据在文件当中的位置
}

// LogRecord 编码，返回字节数组和长度
func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	length := int64(len(logRecord.Key) + len(logRecord.Value) + 2) // 2 for the type byte
	data := make([]byte, length)
	

	data[0] = logRecord.Type
	copy(data[1:], logRecord.Key)
	copy(data[1+len(logRecord.Key):], logRecord.Value)
	
	return data, length
}
