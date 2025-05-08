package fio

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 销毁文件。测试完成后，删除测试文件
func destroyFile(fileName string) {
	if err := os.RemoveAll(fileName); err != nil {
		log.Printf("Failed to remove file %s: %v", fileName, err)
		// panic(error) 测试会不通过
		// 在测试中，避免使用 panic，而应该使用日志记录错误
	}
}

func TestNewFileIOManager(t *testing.T) {
	path := filepath.Join(os.TempDir(), "a.data")
	fio, err := NewFileIOManager(path)
	// 生成的文件路径为：C:\Users\yszgb\AppData\Local\Temp\a.data
	defer destroyFile(path) // defer 延迟执行，在函数返回前执行，确保文件被删除

	assert.Nil(t, err)
	assert.NotNil(t, fio)
}

func TestFileIO_Write(t *testing.T) {
	path := filepath.Join(os.TempDir(), "b.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	n, err := fio.Write([]byte("")) // 写入空字节数组，返回 0，nil
	assert.Equal(t, 0, n)           // 比较写入的字节数是否为 0
	assert.Nil(t, err)              // 比较错误是否为 nil

	n, err = fio.Write([]byte("bitcask kv"))
	// fio.fd.Sync() // 同步数据到磁盘，确保数据持久化。否则，相同的内容不会写入
	// t.Log(n, err) // 打印 10，nil，写入 10 个字节，错误为 nil
	assert.Equal(t, 10, n)
	assert.Nil(t, err)

	n, err = fio.Write([]byte("storage"))
	// t.Log(n, err) // 打印 7，nil，写入 7 个字节，错误为 nil
	assert.Equal(t, 7, n)
	assert.Nil(t, err)
}

func TestFileIO_Read(t *testing.T) {
	path := filepath.Join(os.TempDir(), "0001.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	_, err = fio.Write([]byte("key-a"))
	assert.Nil(t, err)

	_, err = fio.Write([]byte("key-b"))
	assert.Nil(t, err)

	b1 := make([]byte, 5)
	n, err := fio.Read(b1, 0) // 从第 0 个字节开始读取 5 个字节，返回读取的字节数和错误
	t.Log(b1, n)              // 打印 [107 101 121 45 97] 5，读取 5 个字节
	// 为何是 107，101，121，45，97？
	// 因为写入的是 "key-a"，每个字符占用 1 个字节，所以读取的是 "key-a" 的 ASCII 码
	// []byte 默认以整数切片形式存储，每个元素占用 1 个字节
	// 要显示字符串，需要显式指定 t.Log(string(b), n) // 正确输出 "key-a" 5
	assert.Equal(t, 5, n)                // 比较读取的字节数是否为 5
	assert.Nil(t, err)                   // 比较错误是否为 nil
	assert.Equal(t, []byte("key-a"), b1) // 比较读取的字节数组是否为 "key-a"

	b2 := make([]byte, 5)
	_, err = fio.Read(b2, 5)
	t.Log(string(b2), err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("key-b"), b2)
	assert.Nil(t, err)
}

func TestFileIO_Sync(t *testing.T) {
	tests := []struct {
		name    string
		fio     *FileIO
		wantErr bool
	}{
		{
			name: "Sync on valid file",
			fio: func() *FileIO {
				path := filepath.Join(os.TempDir(), "0001.data")
				fio, _ := NewFileIOManager(path)
				return fio
			}(),
			wantErr: false, // 期望没有错误
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer destroyFile(filepath.Join(os.TempDir(), "0001.data"))
			if err := tt.fio.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("FileIO.Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// path := filepath.Join(os.TempDir(), "0001.data")
	// fio, err := NewFileIOManager(path)
	// defer destroyFile(path)
	// assert.Nil(t, err)
	// assert.NotNil(t, fio)

	// err = fio.Sync()
	// assert.Nil(t, err)
}

func TestFileIO_Close(t *testing.T) {
	tests := []struct {
		name    string
		fio     *FileIO
		wantErr bool
	}{
		{
			name: "Close on valid file",
			fio: func() *FileIO {
				path := filepath.Join(os.TempDir(), "0001.data")
				fio, _ := NewFileIOManager(filepath.Join(path))
				return fio
			}(),
			wantErr: false, // 期望没有错误
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer destroyFile(filepath.Join(os.TempDir(), "0001.data"))
			if err := tt.fio.Close(); (err != nil) != tt.wantErr {
				t.Errorf("FileIO.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// path := filepath.Join(os.TempDir(), "0001.data")
	// fio, err := NewFileIOManager(path)
	// defer destroyFile(path)
	// assert.Nil(t, err)
	// assert.NotNil(t, fio)

	// err = fio.Close()
	// assert.Nil(t, err)
}
